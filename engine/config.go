package engine

import (
	conf "github.com/daviddengcn/go-ljson-conf"
)

type Config struct {
	*conf.Conf
}

func LoadConfig(fn string) (*Config, error) {
	cf, err := conf.Load(fn)
	if err != nil {
		return nil, err
	}

	this := new(Config)
	this.Conf = cf

	return this, nil
}
