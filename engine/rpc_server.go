package engine

import (
	"errors"
	"fmt"
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/thrift/lib/go/thrift"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	activeSessionN int64
	cumSessions    int64 // cumulated
	cumCalls       int64
	cumCallErrs    int64
	cumCallSlow    int64 // TODO
	stats          *engineStats

	mutex  sync.Mutex
	errors map[string]int32

	quit       chan bool
	dispatcher *rpcDispatcher

	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory
}

func NewTFunServer(engine *Engine,
	preforkMode bool,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	this := &TFunServer{
		quit:                   make(chan bool),
		stats:                  newEngineStats(),
		errors:                 make(map[string]int32, 1<<10),
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,  // TServerSocket
		inputTransportFactory:  transportFactory, // TBufferedTransportFactory
		outputTransportFactory: transportFactory, // TBufferedTransportFactory
		inputProtocolFactory:   protocolFactory,  // TBinaryProtocolFactory
		outputProtocolFactory:  protocolFactory,  // TBinaryProtocolFactory
	}
	this.dispatcher = newRpcDispatcher(preforkMode,
		config.Engine.Rpc.MaxOutstandingSessions,
		this.handleSession)
	return this
}

// thrift.TServer.Stop
func (this *TFunServer) Stop() error {
	close(this.quit)
	this.serverTransport.Interrupt()
	this.serverTransport.Close() // accept tcp [::]:9001: use of closed network connection
	return nil
}

// thrift.TServer.Serve
func (this *TFunServer) Serve() error {
	if err := this.serverTransport.Listen(); err != nil {
		return err
	}

	// start the stats counter
	go this.stats.Start(time.Now(),
		config.Engine.Rpc.StatsOutputInterval,
		config.Engine.MetricsLogfile)

	if config.Engine.Rpc.StatsOutputInterval > 0 {
		go this.showStats(config.Engine.Rpc.StatsOutputInterval)
	}

	// register to etcd
	// once registered, other peers will connect to me
	// so, must be after Listen ready
	if config.Engine.ServerMode &&
		config.Engine.EtcdSelfAddr != "" {
		etclib.BootService(config.Engine.EtcdSelfAddr, etclib.SERVICE_FAE)

		log.Info("etcd self[%s] registered", config.Engine.EtcdSelfAddr)
	}

	const (
		SLEEP_STEP       = 2
		ACCEPT_MIN_SLEEP = time.Millisecond * 10
		ACCEPT_MAX_SLEEP = ACCEPT_MIN_SLEEP * SLEEP_STEP * SLEEP_STEP *
			SLEEP_STEP * SLEEP_STEP
	)
	var delay = ACCEPT_MIN_SLEEP
	for {
		select {
		case <-this.quit:
			log.Info("RPC server quit...")
			return errors.New("RPC server quit")

		default:
		}

		client, err := this.serverTransport.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Warn("Accept temporary: %s, sleeping %dms", ne.Error(),
					delay/time.Millisecond)
				time.Sleep(delay)
				delay *= SLEEP_STEP
				if delay > ACCEPT_MAX_SLEEP {
					delay = ACCEPT_MAX_SLEEP
				}
			} else {
				log.Error("Accept: %s", err.Error())
			}

			continue
		}

		this.dispatcher.Dispatch(client)
		delay = ACCEPT_MIN_SLEEP
	}

	return nil
}

func (this *TFunServer) Runtime() map[string]interface{} {
	r := make(map[string]interface{})
	r["sessions.active"] = atomic.LoadInt64(&this.activeSessionN)
	r["sessions.all"] = atomic.LoadInt64(&this.cumSessions)
	r["call.all"] = atomic.LoadInt64(&this.cumCalls)
	r["call.err"] = atomic.LoadInt64(&this.cumCallErrs)
	r["call.err.msg"] = this.errors
	r["dispatcher"] = this.dispatcher.Runtime()
	r["qps"] = fmt.Sprintf("1m:%.0f, 5m:%.0f 15m:%.0f avg:%.0f",
		this.stats.CallPerSecond.Rate1(),
		this.stats.CallPerSecond.Rate5(),
		this.stats.CallPerSecond.Rate15(),
		this.stats.CallPerSecond.RateMean())

	return r
}

func (this *TFunServer) showStats(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for _ = range ticker.C {
		log.Info("rpc: {session.on:%d/%s, call.err:%s/%s, qps:{1m:%.0f, 5m:%.0f 15m:%.0f avg:%.0f}}",
			atomic.LoadInt64(&this.activeSessionN),
			gofmt.Comma(atomic.LoadInt64(&this.cumSessions)),
			gofmt.Comma(atomic.LoadInt64(&this.cumCallErrs)),
			gofmt.Comma(atomic.LoadInt64(&this.cumCalls)),
			this.stats.CallPerSecond.Rate1(),
			this.stats.CallPerSecond.Rate5(),
			this.stats.CallPerSecond.Rate15(),
			this.stats.CallPerSecond.RateMean())
	}
}

func (this *TFunServer) handleSession(client thrift.TTransport) {
	var (
		calls           int64 // #calls within this session
		errs            int64 // #errs within this session
		t1              = time.Now()
		currentSessionN = atomic.AddInt64(&this.activeSessionN, 1)
		tcpClient       = client.(*thrift.TSocket).Conn().(*net.TCPConn)
		remoteAddr      = tcpClient.RemoteAddr().String()
		processor       = this.processorFactory.GetProcessor(client)
		inputTransport  = this.inputTransportFactory.GetTransport(client)
		outputTransport = this.outputTransportFactory.GetTransport(client)
		inputProtocol   = this.inputProtocolFactory.GetProtocol(inputTransport)
		outputProtocol  = this.outputProtocolFactory.GetProtocol(outputTransport)
	)
	atomic.AddInt64(&this.cumSessions, 1)
	log.Debug("session[%s]#%d open", remoteAddr, currentSessionN)

	if calls, errs = this.serveCalls(tcpClient, remoteAddr, processor,
		inputProtocol, outputProtocol); errs > 0 {
		atomic.AddInt64(&this.cumCallErrs, errs)
	}
	atomic.AddInt64(&this.cumCalls, calls)

	// server actively closes the socket
	if inputTransport != nil {
		inputTransport.Close()
	}
	if outputTransport != nil {
		outputTransport.Close()
	}

	currentSessionN = atomic.AddInt64(&this.activeSessionN, -1) + 1
	elapsed := time.Since(t1)

	if errs > 0 {
		log.Warn("session[%s]#%d %d calls in %s, errs:%d", remoteAddr,
			currentSessionN, calls, elapsed, errs)
	} else {
		log.Trace("session[%s]#%d %d calls in %s", remoteAddr,
			currentSessionN, calls, elapsed)
	}
}

func (this *TFunServer) serveCalls(tcpClient *net.TCPConn,
	remoteAddr string,
	processor thrift.TProcessor,
	inputProtocol thrift.TProtocol,
	outputProtocol thrift.TProtocol) (callsN int64, errsN int64) {
	var (
		t1      time.Time
		elapsed time.Duration
	)

	for {
		t1 = time.Now()
		if config.Engine.Rpc.IoTimeout > 0 { // read + write
			tcpClient.SetDeadline(t1.Add(config.Engine.Rpc.IoTimeout))
		}

		_, ex := processor.Process(inputProtocol, outputProtocol)
		callsN++ // call num increment first anyway

		elapsed = time.Since(t1)
		this.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		this.stats.CallPerSecond.Mark(1)

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
		if err, isTransportEx := ex.(thrift.TTransportException); isTransportEx {
			if err.TypeId() != thrift.END_OF_FILE {
				// e,g. connection reset by peer: read a socket that was closed by peer
				// e,g. broken pipe: write to socket that was closed by peer
				// e,g. read tcp i/o timeout
				log.Error("transport[%s]: %s", remoteAddr, ex.Error())
				errsN++
				this.saveCallError(err)
			} else {
				// EOF is not err, its normal end of session
				err = nil
			}

			callsN-- // in case of transport err, the call didn't finish
			this.stats.CallPerSession.Update(callsN)

			// for transport err, server always stop the session
			return
		}

		// TProtocolException should never happen
		// so ex MUST be servant generated TApplicationException
		// e,g Error 1064: You have an error in your SQL syntax
		errsN++
		this.saveCallError(ex)

		// it must be TApplicationException
		// the central place to log call err
		// servant needn't dup err log
		log.Error("caller[%s]: %s", remoteAddr, ex.Error())
	}

	this.stats.CallPerSession.Update(callsN)
	return
}

func (this *TFunServer) saveCallError(err error) {
	this.mutex.Lock()
	if len(this.errors) > (20 << 10) { // avoid OOM
		this.errors = make(map[string]int32, 1<<10)
	}
	this.errors[err.Error()]++
	this.mutex.Unlock()
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
