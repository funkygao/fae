package main

var (
	options struct {
		configFile  string
		showVersion bool
		logFile     string
		tick        int
		cpuprof     bool
		memprof     bool
		blockprof   bool
		lockFile    string
		logLevel    string
	}
)

const (
	PROFILER_DIR = "profiler"
	USAGE        = `fae - Fun App Engine

Flags:
`
)
