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

	currentSessionN := atomic.AddInt64(&this.activeSessionN, 1)
	defer atomic.AddInt64(&this.activeSessionN, -1)

	remoteAddr := transport.(*thrift.TSocket).Conn().(*net.TCPConn).RemoteAddr().String()
	if currentSessionN > config.Engine.Rpc.WarnTooManySessionsThreshold {
		log.Warn("session[%s] open, too many sessions: %d",
			remoteAddr, currentSessionN)
	} else {
		log.Debug("session[%s] open", remoteAddr)
	}

	var (
		t1    = time.Now()
		calls int64
		errs  int64
	)
	if calls, errs = this.processRequests(transport); errs > 0 {
		this.engine.svt.AddErr(errs)
	}

	elapsed := time.Since(t1)
	if errs > 0 {
		log.Warn("session[%s] %d calls in %s, errs:%d", remoteAddr, calls, elapsed, errs)
	} else {
		log.Trace("session[%s] %d calls in %s", remoteAddr, calls, elapsed)
	}

}

func (this *TFunServer) processRequests(client thrift.TTransport) (callsN int64, errsN int64) {
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
		remoteAddr   = tcpClient.RemoteAddr().String()
	)

	for {
		t1 = time.Now()
		if rpcIoTimeout > 0 { // read + write
			tcpClient.SetDeadline(t1.Add(rpcIoTimeout))
		}

		_, ex := processor.Process(inputProtocol, outputProtocol)
		callsN++ // call num increment first anyway

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
				// e,g. connection reset by peer
				// e,g. broken pipe
				// e,g. read tcp i/o timeout
				log.Error("transport[%s]: %s", remoteAddr, ex.Error())
				errsN++
			} else {
				// EOF is not err, its normal end of session
				err = nil
			}

			callsN-- // in case of transport err, the call didn't finish
			this.engine.stats.CallPerSession.Update(callsN)

			// for transport err, server always stop the session
			return
		}

		// TProtocolException should never happen
		// so ex MUST be servant generated TApplicationException
		// e,g Error 1064: You have an error in your SQL syntax
		errsN++

		// the central place to log call err
		// servant needn't dup err log
		log.Error("caller[%s]: %s", remoteAddr, ex.Error())

		// Peek: there is more data to be read or the remote side is still open?
		if !inputProtocol.Transport().Peek() {
			break
		}
	}

	this.engine.stats.CallPerSession.Update(callsN)
	return
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
		log.Info("rpc: {active_sessions:%d, qps:{1m:%.1f, 5m:%.1f 15m:%.1f avg:%.1f}}",
			atomic.LoadInt64(&this.activeSessionN),
			this.engine.stats.CallPerSecond.Rate1(),
			this.engine.stats.CallPerSecond.Rate5(),
			this.engine.stats.CallPerSecond.Rate15(),
			this.engine.stats.CallPerSecond.RateMean())
	}
}
