package main

import (
	"github.com/funkygao/fxi/engine"
)

var (
	BuildID = "unknown" // git version id, passed in from shell

	options struct {
		configFile  string
		showVersion bool
		logFile     string
		tick        int
		cpuprof     string
		memprof     string
		lockFile    string
		logLevel    string
	}
)

const (
	USAGE = `fxi - Funplus eXchange Interface

Flags:
`
)
