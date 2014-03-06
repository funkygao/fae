package config

import (
	conf "github.com/funkygao/jsconf"
	"time"
)

type ConfigBreaker struct {
	FailureAllowance uint
	RetryTimeout     time.Duration
}

func (this *ConfigBreaker) loadConfig(cf *conf.Conf) {
	this.FailureAllowance = uint(cf.Int("failure_allowance", 5))
	this.RetryTimeout = time.Second * time.Duration(cf.Int("retry_timeout", 10))
}
