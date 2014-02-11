package servant

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"testing"
)

func TestContextInfo(t *testing.T) {
	ctx := rpc.NewContext()
	ctx.Caller = "POST+/facebook/getPaymentRequestId/+34ca2cf6"

	fun := FunServantImpl{}
	info := fun.contextInfo(ctx)
	assert.Equal(t, true, info.Valid())
	assert.Equal(t, "POST", info.httpMethod)
	assert.Equal(t, "/facebook/getPaymentRequestId/", info.uri)
	assert.Equal(t, "34ca2cf6", info.seqId)

	ctx.Caller = ""
	info = fun.contextInfo(ctx)
	assert.Equal(t, false, info.Valid())
}

func TestNormalizedKind(t *testing.T) {
	fun := FunServantImpl{}
	assert.Equal(t, "log", fun.normalizedKind("database.log"))
	assert.Equal(t, "db5", fun.normalizedKind("db5"))
	assert.Equal(t, "db3", fun.normalizedKind("database.db3"))
}
