package main

import (
	"fmt"
	"github.com/funkygao/fxi/engine"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/signal"
	"net/http"
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

	e := engine.NewEngine(options.configFile).LoadConfigFile()
	e.RegisterHttpApi("/ver", func(w http.ResponseWriter,
		req *http.Request, params map[string]interface{}) (interface{}, error) {
		output := make(map[string]interface{})
		output["ver"] = BuildID
		return output, nil
	})
	e.ServeForever()

	shutdown()
}
