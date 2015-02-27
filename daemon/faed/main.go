package main

import (
	"fmt"
	"github.com/funkygao/fae/engine"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/profile"
	"github.com/funkygao/golib/server"
	"github.com/funkygao/golib/signal"
	log "github.com/funkygao/log4go"
	_log "log"
	"os"
	"runtime/debug"
	"syscall"
	"time"
)

var (
	engineRunner *engine.Engine
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	server.SetupLogging(options.logFile, options.logLevel,
		options.crashLogFile, options.alarmLogSock, options.alarmTag)

	// thrift lib use "log", so we also need to customize its behavior
	_log.SetFlags(_log.Ldate | _log.Ltime | _log.Lshortfile)

	if options.kill {
		s := server.NewServer("fae")
		s.LoadConfig(options.configFile)
		s.Launch()

		// stop new requests
		engine.NewEngine().
			LoadConfig(options.configFile, s.Conf).
			UnregisterEtcd()

		// finish all outstanding RPC sessions
		if err := server.SignalProcess(options.lockFile, syscall.SIGUSR1); err != nil {
			fmt.Fprintf(os.Stderr, "stop failed: %s\n", err)
		}

		cleanup() // TODO wait till that faed process terminates, who will do the cleanup

		fmt.Println("faed killed")

		os.Exit(0)
	}

	if options.lockFile != "" {
		if locking.InstanceLocked(options.lockFile) {
			fmt.Fprintf(os.Stderr, "Another instance is running, exit...\n")
			os.Exit(1)
		}

		locking.LockInstance(options.lockFile)
	}

}

func main() {
	defer func() {
		cleanup()

		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	if options.cpuprof || options.memprof {
		cf := &profile.Config{
			Quiet:        true,
			ProfilePath:  "prof",
			CPUProfile:   options.cpuprof,
			MemProfile:   options.memprof,
			BlockProfile: options.blockprof,
		}

		defer profile.Start(cf).Stop()
	}

	log.Info("%s", `
     ____      __      ____ 
    ( ___)    /__\    ( ___)
     )__)    /(__)\    )__) 
    (__)    (__)(__)  (____)`)

	s := server.NewServer("fae")
	s.LoadConfig(options.configFile)
	s.Launch()

	go server.RunSysStats(time.Now(), options.tick)

	engineRunner = engine.NewEngine()
	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		shutdown()
		engineRunner.StopRpcServe()
	})

	engineRunner.LoadConfig(options.configFile, s.Conf).
		ServeForever()
}
