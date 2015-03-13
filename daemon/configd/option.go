package main

import (
	"flag"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"io/ioutil"
)

var (
	options struct {
		configFile    string
		showVersion   bool
		logFile       string
		logLevel      string
		kill          bool
		crashLogFile  string
		statsInterval int
	}

	config struct {
		etcServers           []string
		faeTemplateFile      string
		faeTargetFile        string
		actorTemplateFile    string
		actorTargetFile      string
		maintainTemplateFile string
		maintainTargetFile   string
	}
)

var (
	faeTemplateContents      string
	actorTemplateContents    string
	maintainTemplateContents string
)

func parseFlags() {
	flag.BoolVar(&options.kill, "k", false, "kill")
	flag.StringVar(&options.configFile, "conf", "etc/configd.cf", "config file")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.IntVar(&options.statsInterval, "interval", 10, "show stats per how many minutes")
	flag.StringVar(&options.crashLogFile, "crashlog", "panic.dump", "crash log file")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")

	flag.Parse()
}

func loadConfig(cf *conf.Conf) {
	config.etcServers = cf.StringList("etcd_servers", nil)
	config.faeTemplateFile = cf.String("fae_template_file", "")
	config.faeTargetFile = cf.String("fae_target_file", "")
	config.maintainTemplateFile = cf.String("maintain_template_file", "")
	config.maintainTargetFile = cf.String("maintain_target_file", "")

	log.Debug("config: %+v", config)
}

func loadTemplates() {
	if config.faeTemplateFile != "" {
		body, err := ioutil.ReadFile(config.faeTemplateFile)
		if err != nil {
			log.Error("template[%s]: %s", config.faeTemplateFile, err)
		} else {
			faeTemplateContents = string(body)

			log.Info("template[%s] loaded", config.faeTemplateFile)
		}
	}

	if config.actorTemplateFile != "" {
		body, err := ioutil.ReadFile(config.actorTemplateFile)
		if err != nil {
			log.Error("template[%s]: %s", config.actorTemplateFile, err)
		} else {
			maintainTemplateContents = string(body)

			log.Info("template[%s] loaded", config.actorTemplateFile)
		}
	}

	if config.maintainTemplateFile != "" {
		body, err := ioutil.ReadFile(config.maintainTemplateFile)
		if err != nil {
			log.Error("template[%s]: %s", config.maintainTemplateFile, err)
		} else {
			maintainTemplateContents = string(body)

			log.Info("template[%s] loaded", config.maintainTemplateFile)
		}
	}
}
