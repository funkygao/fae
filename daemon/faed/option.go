package main

import (
	"flag"
	"fmt"
	log "github.com/funkygao/log4go"
	_log "log"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func parseFlags() {
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.configFile, "conf", "etc/faed.cf", "config file")
	flag.StringVar(&options.lockFile, "lockfile", "faed.lock", "lockfile path")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")
	flag.IntVar(&options.tick, "tick", 60*10, "watchdog ticker length in seconds")
	flag.StringVar(&options.cpuprof, "cpuprof", "", "cpu profiling file")
	flag.StringVar(&options.memprof, "memprof", "", "memory profiling file")
	flag.Usage = showUsage

	flag.Parse()

	if options.tick <= 0 {
		panic("tick must be possitive")
	}
}

func showUsage() {
	fmt.Fprint(os.Stderr, USAGE)
	flag.PrintDefaults()
}

func setupProfiler() {
	if options.cpuprof != "" {
		f, err := os.Create(options.cpuprof)
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(f)

		log.Info("CPU profiler [%s] enabled", options.cpuprof)
	}

	if options.memprof != "" {
		log.Info("MEM profiler [%s] enabled", options.memprof)
	}
}

func setupLogging(loggingLevel, logFile string) {
	level := log.DEBUG
	switch loggingLevel {
	case "info":
		level = log.INFO
	case "warn":
		level = log.WARNING
	case "error":
		level = log.ERROR
	}

	for _, filter := range log.Global {
		filter.Level = level
	}

	if logFile == "stdout" {
		log.AddFilter("stdout", level, log.NewConsoleLogWriter())
	} else {
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0744); err != nil {
			panic(err)
		}

		writer := log.NewFileLogWriter(logFile, false)
		log.AddFilter("file", level, writer)
		writer.SetFormat("[%d %T] [%L] (%S) %M")
		writer.SetRotate(true)
		writer.SetRotateSize(0)
		writer.SetRotateLines(0)
		writer.SetRotateDaily(true)
	}

	// thrift lib use "log", so we also need to customize its behavior
	_log.SetFlags(_log.Ldate | _log.Ltime | _log.Lshortfile)
}
