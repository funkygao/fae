package servant

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"testing"
)

func TestContextCallerInfo(t *testing.T) {
	ctx := rpc.NewContext()
	ctx.Caller = "POST+/facebook/getPaymentRequestId/+34ca2cf6"

	fun := FunServantImpl{}
	info := fun.callerInfo(ctx)
	assert.Equal(t, true, info.Valid())
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, "/facebook/getPaymentRequestId/", info.uri)
	assert.Equal(t, "34ca2cf6", info.seqId)

	ctx.Caller = ""
	info = fun.callerInfo(ctx)
	assert.Equal(t, false, info.Valid())
}
