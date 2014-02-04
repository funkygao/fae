package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	conf "github.com/funkygao/jsconf"
	"testing"
)

func setupServant() *FunServantImpl {
	cf, _ := conf.Load("../etc/faed.cf")
	section, _ := cf.Section("servants")
	config.LoadServants(section)
	return NewFunServant(config.Servants)
}

func BenchmarkMcSet(b *testing.B) {
	servant := setupServant()
	b.ReportAllocs()

	ctx := rpc.NewReqCtx()
	ctx.Caller = "me"
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "foo", []byte("bar"), 0)
	}
}
