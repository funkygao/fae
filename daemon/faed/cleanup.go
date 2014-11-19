package main

import (
	"fmt"
	"github.com/funkygao/golib/locking"
	log "github.com/funkygao/log4go"
	"io/ioutil"
	"os"
	"strconv"
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

func killFaed() {
	filebody, err := ioutil.ReadFile(options.lockFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to die[%s]: %s", options.lockFile, err)
		return
	}

	pid, err := strconv.Atoi(string(filebody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to die[%s]: %s", options.lockFile, err)
		return
	}

	faedProcess, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to die[%d]: %s", pid, err)
		return
	}

	faedProcess.Kill()

	locking.UnlockInstance(options.lockFile)
}
