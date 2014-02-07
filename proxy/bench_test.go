package proxy

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"testing"
)

func BenchmarkServantPing(b *testing.B) {
	servant, closer := Servant(":9001")
	defer closer.Close()

	ctx := rpc.NewContext()
	ctx.Caller = "me"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.Ping(ctx)

	}
	b.SetBytes(10)
}

func BenchmarkServantMcSet(b *testing.B) {
	servant, closer := Servant(":9001")
	defer closer.Close()

	ctx := rpc.NewContext()
	ctx.Caller = "me"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "foo", []byte("bar"), 0)
	}
	b.SetBytes(10)
}
