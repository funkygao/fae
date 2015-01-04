package config

import (
	conf "github.com/funkygao/jsconf"
	"os"
)

func LoadEngineConfig(configFile string, cf *conf.Conf) {
	Engine = new(ConfigEngine)
	Engine.configFile = configFile
	var err error
	Engine.configFileLastStat, err = os.Stat(configFile)
	if err != nil {
		panic(err)
	}
	Engine.LoadConfig(cf)

	go Engine.runWatchdog()
}
