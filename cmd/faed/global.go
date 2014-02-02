package main

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
	USAGE = `fae - Funplus App Engine

Flags:
`
)
