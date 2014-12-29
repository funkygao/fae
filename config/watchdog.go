package config

import (
	log "github.com/funkygao/log4go"
	"os"
	"time"
)

func (this *ConfigEngine) runWatchdog() {
	ticker := time.NewTicker(this.ReloadWatchdogInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		stat, _ := os.Stat(Engine.configFile)
		if stat.ModTime() != Engine.configFileLastStat.ModTime() {
			Engine.configFileLastStat = stat

			// TODO
			log.Info("config[%s] reloaded", Engine.configFile)

		}
	}

}
