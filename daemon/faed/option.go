package main

import (
	"flag"
)

var (
	options struct {
		configFile  string
		showVersion bool
		logFile     string
		logLevel    string
		tick        int
		cpuprof     bool
		memprof     bool
		blockprof   bool
		lockFile    string
	}
)

func parseFlags() {
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.configFile, "conf", "etc/faed.cf", "config file")
	flag.StringVar(&options.lockFile, "lockfile", "faed.lock", "lockfile path")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")
	flag.IntVar(&options.tick, "tick", 60*10, "watchdog ticker length in seconds")
	flag.BoolVar(&options.cpuprof, "cpuprof", false, "enable cpu profiling")
	flag.BoolVar(&options.memprof, "memprof", false, "enable memory profiling")
	flag.BoolVar(&options.blockprof, "blockprof", false, "enable block profiling")

	flag.Parse()

	if options.tick <= 0 {
		panic("tick must be possitive")
	}
}
