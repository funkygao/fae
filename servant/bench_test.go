package servant

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	rand_ "github.com/funkygao/golib/rand"
	"github.com/funkygao/golib/server"
	conf "github.com/funkygao/jsconf"
	"strings"
	"testing"
	"time"
)

func setupServant() *FunServantImpl {
	server.SetupLogging(".canbedeleted.test.log", "info", "", "", "")

	cf, _ := conf.Load("../etc/faed.cf")
	config.LoadEngineConfig(cf)
	server.LaunchHttpServer(":9999", "")
	return NewFunServant(config.Engine.Servants)
}

// 880 ns/op
func BenchmarkIsSelectQuery(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(strings.ToLower(sql), "select")
	}
}

// 64.0 ns/op
func BenchmarkOptimizedIsSelectQuery(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(strings.ToLower(sql[:len("select")]), "select")
	}
}

// 9 ns/op
func BenchmarkIsSelectQueryWithoutLowcase(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(sql, "select")
	}
}

func BenchmarkSizedString12(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rand_.SizedString(12)
	}
}

func BenchmarkGetSession(b *testing.B) {
	servant := setupServant()
	b.ReportAllocs()
	ctx := rpc.NewContext()
	ctx.Reason = "Map:enterKingdomBlock"
	for i := 0; i < b.N; i++ {
		ctx.Rid = time.Now().UnixNano()
		servant.getSession(ctx)
	}
}

// 683 ns/op
func BenchmarkSprintfSql(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("SELECT %s FROM %s WHERE %s", "*", "UserInfo", "uid=?")
	}

}

// 166 ns/op
func BenchmarkRawConcatSql(b *testing.B) {
	b.ReportAllocs()
	table := "UserInfo"
	column := "*"
	where := "uid=?"
	for i := 0; i < b.N; i++ {
		_ = "SELECT " + column + " FROM " + table + " WHERE " + where
	}
}

func BenchmarkMcSet(b *testing.B) {
	servant := setupServant()
	b.ReportAllocs()

	ctx := rpc.NewContext()
	ctx.Reason = "benchmark"
	ctx.Rid = 1
	data := rpc.NewTMemcacheData()
	data.Data = []byte("bar")
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "default", "foo", data, 0)
	}
}

func BenchmarkPingOnLocalhost(b *testing.B) {
	b.ReportAllocs()

	cf := &config.ConfigProxy{PoolCapacity: 1}
	client, err := proxy.New(cf).ServantByAddr("localhost:9001")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = "benchmark"
	ctx.Rid = 2

	for i := 0; i < b.N; i++ {
		client.Ping(ctx)
	}
}
