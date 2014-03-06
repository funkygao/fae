package engine

import (
	"errors"
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/funkygao/log4go"
	"net"
	"strings"
	"sync"
	"time"
)

// thrift.TServer implementation
type TFunServer struct {
	stopped bool

	engine                 *Engine
	processorFactory       thrift.TProcessorFactory
	serverTransport        thrift.TServerTransport
	inputTransportFactory  thrift.TTransportFactory
	outputTransportFactory thrift.TTransportFactory
	inputProtocolFactory   thrift.TProtocolFactory
	outputProtocolFactory  thrift.TProtocolFactory

	pool *rpcThreadPool

	mu                  sync.Mutex
	clientConcurrencies map[string]int
}

func NewTFunServer(engine *Engine,
	processor thrift.TProcessor,
	serverTransport thrift.TServerTransport,
	transportFactory thrift.TTransportFactory,
	protocolFactory thrift.TProtocolFactory) *TFunServer {
	this := &TFunServer{
		engine:                 engine,
		processorFactory:       thrift.NewTProcessorFactory(processor),
		serverTransport:        serverTransport,
		inputTransportFactory:  transportFactory,
		outputTransportFactory: transportFactory,
		inputProtocolFactory:   protocolFactory,
		outputProtocolFactory:  protocolFactory,
		clientConcurrencies:    make(map[string]int),
	}
	this.pool = newRpcThreadPool(this.engine.conf.rpc.pm, this.handleClient)
	engine.rpcThreadPool = this.pool

	// start the thread pool
	this.pool.Start()

	// any web frontend got stuck?
	go this.monitorClients()

	return this
}

func (this *TFunServer) Serve() error {
	this.stopped = false
	err := this.serverTransport.Listen()
	if err != nil {
		return err
	}

	for !this.stopped {
		client, err := this.serverTransport.Accept()
		if client != nil {
			this.pool.Dispatch(client)
		}

		if err != nil {
			log.Error("Accept: %v", err)
		}
	}

	return errors.New("rpc server stopped")
}

func (this *TFunServer) handleClient(client interface{}) {
	defer this.engine.stats.CurrentSessions.Dec(1)

	transport, ok := client.(thrift.TTransport)
	if !ok {
		log.Error("Invalid client: %#v", client)
		return
	}

	this.engine.stats.SessionPerSecond.Mark(1)
	this.engine.stats.CurrentSessions.Inc(1)

	if tcpClient, ok := transport.(*thrift.TSocket).Conn().(*net.TCPConn); ok {
		if !this.engine.conf.rpc.tcpNoDelay {
			// golang is tcp no delay by default
			tcpClient.SetNoDelay(false)
		}

		if this.engine.conf.rpc.debugSession {
			log.Debug("Accepted session peer{%s}", tcpClient.RemoteAddr())
		}

		// store client concurrent connections count
		this.mu.Lock()
		p := strings.SplitN(tcpClient.RemoteAddr().String(), ":", 2)
		if len(p) == 2 && p[0] != "" {
			this.clientConcurrencies[p[0]] += 1
			defer func() {
				this.mu.Lock()
				this.clientConcurrencies[p[0]] -= 1
				this.mu.Unlock()
			}()
		}
		this.mu.Unlock()
	}

	this.processSession(transport)
}

// if concurrent conns from same client is too high, it means
// web frontend(php-fpm) got stuck, keep forking children
// TODO
func (this *TFunServer) monitorClients() {
	for {
		for clientIp, concurrentConns := range this.clientConcurrencies {
			if concurrentConns > 200 {
				log.Warn("Client[%s] may got stuck: %d", clientIp, concurrentConns)
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func (this *TFunServer) processSession(client thrift.TTransport) {
	t1 := time.Now()
	remoteAddr := client.(*thrift.TSocket).Conn().RemoteAddr().String()
	if err := this.processRequest(client); err != nil {
		this.engine.stats.TotalFailedSessions.Inc(1)
		log.Error("Session peer{%s}: %s", remoteAddr, err)
	}

	elapsed := time.Since(t1)
	this.engine.stats.SessionLatencies.Update(elapsed.Nanoseconds() / 1e6)
	if this.engine.conf.rpc.debugSession {
		log.Debug("Closed session peer{%s} after %s", remoteAddr, elapsed)
	} else if elapsed.Seconds() > this.engine.conf.rpc.sessionSlowThreshold {
		// slow session
		this.engine.stats.TotalSlowSessions.Inc(1)
		log.Warn("SLOW=%s session peer{%s}", elapsed, remoteAddr)
	}

}

func (this *TFunServer) processRequest(client thrift.TTransport) error {
	processor := this.processorFactory.GetProcessor(client)
	inputTransport := this.inputTransportFactory.GetTransport(client)
	outputTransport := this.outputTransportFactory.GetTransport(client)
	inputProtocol := this.inputProtocolFactory.GetProtocol(inputTransport)
	outputProtocol := this.outputProtocolFactory.GetProtocol(outputTransport)
	if inputTransport != nil {
		defer inputTransport.Close()
	}
	if outputTransport != nil {
		defer outputTransport.Close()
	}

	var (
		t1         time.Time
		elapsed    time.Duration
		remoteAddr = client.(*thrift.TSocket).Conn().RemoteAddr().String()
	)
	for {
		t1 = time.Now()
		ok, err := processor.Process(inputProtocol, outputProtocol)

		elapsed = time.Since(t1)
		this.engine.stats.CallPerSecond.Mark(1)
		this.engine.stats.CallLatencies.Update(elapsed.Nanoseconds() / 1e6)
		if elapsed.Seconds() > this.engine.conf.rpc.callSlowThreshold {
			// slow call
			this.engine.stats.TotalSlowCalls.Inc(1)
			log.Warn("SLOW call=%.3fs, peer{%s}", elapsed.Seconds(), remoteAddr)
		}

		// check transport error
		if err, ok := err.(thrift.TTransportException); ok &&
			err.TypeId() == thrift.END_OF_FILE {
			// remote client closed transport
			return nil
		} else if err != nil {
			// non-EOF transport err
			// e,g. connection reset by peer
			// e,g. broken pipe
			this.engine.stats.TotalFailedCalls.Inc(1)
			return err
		}

		// it is servant generated TApplicationException
		if err != nil {
			this.engine.stats.TotalFailedCalls.Inc(1)
			log.Error("Servant call peer{%s}: %s", remoteAddr, err)
		}

		if !ok || !inputProtocol.Transport().Peek() {
			break
		}
	}

	return nil
}

func (this *TFunServer) Stop() error {
	this.stopped = true
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
