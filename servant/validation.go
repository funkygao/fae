package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func validateContext(ctx *rpc.Context) error {
	if ctx.Rid == "" || ctx.Reason == "" {
		log.Error("Invalid context: %s", ctx.String())
		return ErrInvalidContext
	}

	return nil
}
