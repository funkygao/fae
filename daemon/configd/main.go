package main

import (
	"fmt"
	"github.com/funkygao/golib/server"
	"runtime/debug"
	"time"
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile, "", "")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	s := server.NewServer("configd")
	s.LoadConfig(options.configFile)
	s.Launch()

	go server.RunSysStats(s.StartedAt,
		time.Minute*time.Duration(options.statsInterval))

	monitorForever(s.Conf)
}
