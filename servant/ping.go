package servant

import (
	log "code.google.com/p/log4go"
)

func (this *FunServantImpl) Ping() (r string, err error) {
	log.Debug("ping")
	return "pong", nil
}
