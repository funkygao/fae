package servant

import (
	"encoding/json"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	conf "github.com/funkygao/jsconf"
	"labix.org/v2/mgo/bson"
	"strings"
	"testing"
)

func setupServant() *FunServantImpl {
	cf, _ := conf.Load("../etc/faed.cf")
	section, _ := cf.Section("servants")
	config.LoadServants(section)
	server.LaunchHttpServ(":9999", "")
	return NewFunServant(config.Servants)
}

func BenchmarkMcSet(b *testing.B) {
	servant := setupServant()
	b.ReportAllocs()

	ctx := rpc.NewContext()
	ctx.Caller = "me"
	data := rpc.NewTMemcacheData()
	data.Data = []byte("bar")
	for i := 0; i < b.N; i++ {
		servant.McSet(ctx, "default", "foo", data, 0)
	}
}

func BenchmarkGoMap(b *testing.B) {
	b.ReportAllocs()
	var x map[int]bool = make(map[int]bool)
	for i := 0; i < b.N; i++ {
		x[i] = true
	}
}

func BenchmarkJson(b *testing.B) {
	m := bson.M{
		"userId": 343434,
		"gendar": "F",
		"info": bson.M{
			"city":    "beijing",
			"hobbies": []string{"a", "b"}}}
	for i := 0; i < b.N; i++ {
		json.Marshal(m)
	}
}

func BenchmarkBson(b *testing.B) {
	m := bson.M{
		"userId": 343434,
		"gendar": "F",
		"info": bson.M{
			"city":    "beijing",
			"hobbies": []string{"a", "b"}}}
	for i := 0; i < b.N; i++ {
		bson.Marshal(m)
	}
}

// 880 ns/op
func BenchmarkIsSelectQuery(b *testing.B) {
	b.ReportAllocs()
	var sql = "select * from UserInfo where uid=? and power>?"
	for i := 0; i < b.N; i++ {
		strings.HasPrefix(strings.ToLower(sql), "select")
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
