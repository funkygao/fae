package main

import (
	"github.com/funkygao/golib/locking"
	log "github.com/funkygao/log4go"
	"os"
)

func cleanup() {
	if options.lockFile != "" {
		locking.UnlockInstance(options.lockFile)
		log.Debug("Cleanup lock %s", options.lockFile)
	}
}

func shutdown() {
	cleanup()

	log.Info("Terminated")

	os.Exit(0)
}
