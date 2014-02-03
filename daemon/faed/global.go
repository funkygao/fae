package main

var (
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
	USAGE = `fae - Fun App Engine

Flags:
`
)
