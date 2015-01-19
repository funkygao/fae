package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"net"
	"sync/atomic"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	quit           chan bool
	activeSessionN int64

	engine                 *Engine
	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory

	pool *rpcThreadPool
}

func NewTFunServer(engine *Engine,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	this := &TFunServer{
		quit:                   make(chan bool),
		engine:                 engine,
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,
		inputTransportFactory:  transportFactory,
		outputTransportFactory: transportFactory,
		inputProtocolFactory:   protocolFactory,
		outputProtocolFactory:  protocolFactory,
	}
	this.pool = newRpcThreadPool(
		config.Engine.Rpc.MaxOutstandingSessions,
		this.handleSession)
	engine.rpcThreadPool = this.pool

	return this
}

func (this *TFunServer) Serve() error {
	const stoppedError = "RPC server stopped"

	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	if config.Engine.Rpc.StatsOutputInterval > 0 {
		go this.showStats(config.Engine.Rpc.StatsOutputInterval)
	}

	// register to etcd
	// once registered, other peers will connect to me
	// so, must be after Listen ready
	if config.Engine.EtcdSelfAddr != "" {
		etclib.BootService(config.Engine.EtcdSelfAddr, etclib.SERVICE_FAE)

		log.Info("etcd self[%s] registered", config.Engine.EtcdSelfAddr)
	}

	for {
		select {
		case <-this.quit:
			// FIXME new conn will timeout, instead of conn close
			log.Info("RPC server quit...")
			return errors.New(stoppedError)

		default:
		}

		client, err := this.serverTransport.Accept()
		if err != nil {
			log.Error("Accept: %v", err)
		} else {
			this.pool.Dispatch(client)
		}
	}

	return errors.New(stoppedError)
}

func (this *TFunServer) handleSession(client interface{}) {
	transport, ok := client.(thrift.TTransport)
	if !ok {
		// should never happen
		log.Error("Invalid client: %#v", client)
		return
	}

	this.engine.stats.SessionPerSecond.Mark(1)
	currentSessionN := atomic.AddInt64(&this.activeSessionN, 1)
	defer atomic.AddInt64(&this.activeSessionN, -1)

	if tcpClient, ok := transport.(*thrift.TSocket).Conn().(*net.TCPConn); ok {
		if currentSessionN > config.Engine.Rpc.WarnTooManySessionsThreshold {
			log.Warn("session[%s] open, too many sessions: %d",
				tcpClient.RemoteAddr(), currentSessionN)
		} else {
			log.Debug("session[%s] open", tcpClient.RemoteAddr())
		}
	} else {
		log.Error("non tcp conn found, should NEVER happen")
		return
	}

	t1 := time.Now()
	remoteAddr := transport.(*thrift.TSocket).Conn().RemoteAddr().String()
	var (
		calls int64
		err   error
	)
	if calls, err = this.processRequests(transport); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	if err != nil {
		log.Error("session[%s] %d calls in %s: %v", remoteAddr, calls, elapsed, err)
	} else {
		log.Trace("session[%s] %d calls in %s", remoteAddr, calls, elapsed)
	}

	if config.Engine.Rpc.SessionSlowThreshold.Seconds() > 0 &&
		elapsed > config.Engine.Rpc.SessionSlowThreshold {
		this.engine.stats.TotalSlowSessions.Inc(1)
	}
}

// FIXME if error occurs, fae will actively close this session: halt on first exception
func (this *TFunServer) processRequests(client thrift.TTransport) (int64, error) {
	processor := this.processorFactory.GetProcessor(client)
	inputTransport := this.inputTransportFactory.GetTransport(client)
	outputTransport := this.outputTransportFactory.GetTransport(client)
	inputProtocol := this.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := this.outputProtocolFactory.GetProtocol(outputTransport)
	defer func() {
		if inputTransport != nil {
			inputTransport.Close()
		}
		if outputTransport != nil {
			outputTransport.Close()
		}
	}()

	var (
		rpcIoTimeout = config.Engine.Rpc.IoTimeout
		t1           time.Time
		elapsed      time.Duration
		tcpClient    = client.(*thrift.TSocket).Conn().(*net.TCPConn)
		callsN       int64
		lastErr      error
	)

	for {
		t1 = time.Now()
		if rpcIoTimeout > 0 { // read + write
			tcpClient.SetDeadline(t1.Add(rpcIoTimeout))
		}

		_, ex := processor.Process(inputProtocol, outputProtocol)
		callsN++

		elapsed = time.Since(t1)
		this.engine.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		this.engine.stats.CallPerSecond.Mark(1)

		if ex == nil {
			// rpc func called/Processed without any error
			continue
		}

		// exception thrown, maybe system wise or app wise

		/*
			thrift exceptions

			TException
				|
				|- TApplicationException
				|- TProtocolException (BAD_VERSION), it should never be thrown, we skip it
				|- TTransportException
		*/
		err, isTransportEx := ex.(thrift.TTransportException)
		if isTransportEx {
			if err.TypeId() != thrift.END_OF_FILE {
				// END_OF_FILE: remote client closed transport, this is normal end of session

				// non-EOF transport err
				// e,g. connection reset by peer
				// e,g. broken pipe
				// e,g. read tcp i/o timeout
				this.engine.stats.TotalFailedCalls.Inc(1)
			}

			this.engine.stats.CallPerSession.Update(callsN)

			// for transport err, server always stop the session
			return callsN, err
		}

		// TProtocolException should never happen
		// so ex MUST be servant generated TApplicationException
		// e,g Error 1064: You have an error in your SQL syntax
		this.engine.stats.TotalFailedCalls.Inc(1)
		lastErr = ex // remember the latest app err

		// Peek: there is more data to be read or the remote side is still open?
		if !inputProtocol.Transport().Peek() {
			break
		}
	}

	this.engine.stats.CallPerSession.Update(callsN)
	return callsN, lastErr
}

func (this *TFunServer) Stop() error {
	close(this.quit)
	this.serverTransport.Interrupt()
	return nil
}

func (this *TFunServer) ProcessorFactory() thrift.TProcessorFactory {
	return this.processorFactory
}

func (this *TFunServer) ServerTransport() thrift.TServerTransport {
	return this.serverTransport
}

func (this *TFunServer) InputTransportFactory() thrift.TTransportFactory {
	return this.inputTransportFactory
}

func (this *TFunServer) OutputTransportFactory() thrift.TTransportFactory {
	return this.outputTransportFactory
}

func (this *TFunServer) InputProtocolFactory() thrift.TProtocolFactory {
	return this.inputProtocolFactory
}

func (this *TFunServer) OutputProtocolFactory() thrift.TProtocolFactory {
	return this.outputProtocolFactory
}

func (this *TFunServer) showStats(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for _ = range ticker.C {
		log.Info("rpc: {active_sessions: %d}", atomic.LoadInt64(&this.activeSessionN))
	}
}
