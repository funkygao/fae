package main

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/golib/locking"
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
