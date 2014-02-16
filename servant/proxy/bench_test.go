package proxy

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"testing"
)

func BenchmarkServantPing(b *testing.B) {
	servant := Servant(":9001")
	defer servant.Transport.Close()

	ctx := rpc.NewContext()
	ctx.Caller = "bench"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.Ping(ctx)

	}
	b.SetBytes(10)
}

func BenchmarkServantMcSet(b *testing.B) {
	servant := Servant(":9001")
	defer servant.Transport.Close()

	ctx := rpc.NewContext()
	ctx.Caller = "bench"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "foo", []byte("bar"), 0)
	}
	b.SetBytes(10)
}
