package config

import (
	conf "github.com/funkygao/jsconf"
)

func LoadEngineConfig(cf *conf.Conf) {
	Engine = new(ConfigEngine)
	Engine.ReloadedChan = make(chan ConfigEngine, 5)
	Engine.LoadConfig(cf)

	go Engine.watchReload()
}
