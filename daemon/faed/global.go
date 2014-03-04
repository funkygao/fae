package main

var (
	options struct {
		configFile  string
		showVersion bool
		logFile     string
		tick        int
		cpuprof     bool
		memprof     bool
		lockFile    string
		logLevel    string
	}
)

const (
	USAGE = `fae - Fun App Engine

Flags:
`
)
