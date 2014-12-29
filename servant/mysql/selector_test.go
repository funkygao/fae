package mysql

import (
	"fmt"
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/str"
	"strings"
	"testing"
)

func TestSelectorStandardEndsWithDigit(t *testing.T) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	assert.Equal(t, true, s.endsWithDigit("AllianceShard8"))
	assert.Equal(t, false, s.endsWithDigit("ShardLookup"))
}

func BenchmarkEndsWithDigit(b *testing.B) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	for i := 0; i < b.N; i++ {
		s.endsWithDigit("UserShard1")
	}
}

// 307 ns/op
func BenchmarkStringsJoin(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	a := []string{
		"SELECT shardId FROM",
		table,
		"WHERE entityId=?",
	}
	for i := 0; i < b.N; i++ {
		strings.Join(a, " ")
	}
}

// 184 ns/op
func BenchmarkStringBuilderConcat(b *testing.B) {
	const (
		s1 = "SELECT shardId FROM "
		s2 = " WHERE entityId=?"
	)

	b.ReportAllocs()
	table := "UserLookup"
	sb := str.NewStringBuilder()
	for i := 0; i < b.N; i++ {
		sb.WriteString(s1)
		sb.WriteString(table)
		sb.WriteString(s2)
		sb.String()
		sb.Reset()
	}
}

// 168 ns/op
func BenchmarkStringConcat(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	for i := 0; i < b.N; i++ {
		_ = "SELECT shardId FROM " +
			table +
			" WHERE entityId=?"
	}
}

// 424 ns/op
func BenchmarkFmtStringConcat(b *testing.B) {
	b.ReportAllocs()
	table := "UserLookup"
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("SELECT shardId FROM %s WHERE entityId=?", table)
	}
}
