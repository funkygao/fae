package servant

import (
	"github.com/funkygao/golib/trace"
	log "github.com/funkygao/log4go"
)

type mongoProtocolLogger struct {
}

func (this *mongoProtocolLogger) Output(calldepth int, s string) error {
	log.Debug("(%s) %s", trace.CallerFuncName(calldepth), s)
	return nil
}
