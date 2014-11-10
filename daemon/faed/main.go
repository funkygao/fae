package main

import (
	"fmt"
	"github.com/funkygao/fae/engine"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/profile"
	"github.com/funkygao/golib/server"
	"github.com/funkygao/golib/signal"
	"os"
	"runtime/debug"
	"syscall"
	"time"
)

func init() {
	parseFlags()

	if options.showVersion {
		engine.ShowVersionAndExit()
	}

	if options.lockFile != "" {
		if locking.InstanceLocked(options.lockFile) {
			fmt.Fprintf(os.Stderr, "Another instance is running, exit...\n")
			os.Exit(1)
		}
		locking.LockInstance(options.lockFile)
	}

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		shutdown()
	})

	setupLogging()
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

	s := server.NewServer("fae")
	s.LoadConfig(options.configFile)
	s.Launch()

	go server.RunSysStats(time.Now(), time.Duration(options.tick)*time.Second)

	engine.NewEngine().
		LoadConfig(s.Conf).
		ServeForever()
}
