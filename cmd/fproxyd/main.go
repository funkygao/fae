package main

import (
	"fmt"
	"github.com/funkygao/fxi/engine"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/signal"
	"os"
	"runtime/debug"
	"syscall"
	"time"
)

func init() {
	parseFlags()

	if options.showVersion {
		showVersionAndExit()
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
}

func main() {
	defer func() {
		cleanup()

		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	setupLogging(options.logLevel, options.logFile)
	setupProfiler()

	ticker := time.NewTicker(time.Second * time.Duration(options.tick))
	go runWatchdog(ticker)
	defer ticker.Stop()

	engine.NewEngine().
		LoadConfigFile(options.configFile).
		ServeForever()

	shutdown()
}
