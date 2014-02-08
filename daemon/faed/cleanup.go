package main

import (
	"github.com/funkygao/golib/locking"
	log "github.com/funkygao/log4go"
	"os"
	"runtime/pprof"
)

func cleanup() {
	if options.lockFile != "" {
		locking.UnlockInstance(options.lockFile)
		log.Debug("Cleanup lock %s", options.lockFile)
	}

	if options.cpuprof != "" {
		pprof.StopCPUProfile()
	}

	if options.memprof != "" {
		f, err := os.Create(options.memprof)
		if err != nil {
			panic(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

func shutdown() {
	cleanup()

	log.Info("Terminated")

	os.Exit(0)
}
