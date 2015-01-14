package mysql

import (
	"errors"
	"github.com/funkygao/assert"
	"testing"
)

func TestIsSystemError(t *testing.T) {
	m := &mysql{}
	assert.NotEqual(t, nil, t)
	err := errors.New("Error Connection failed")
	assert.Equal(t, true, m.isSystemError(err))
	err = errors.New("Error 1054: Unknown column 'curve_internal_id' in 'field list'")
	t.Logf("prefix(%s)", err.Error()[6:])
	assert.Equal(t, false, m.isSystemError(err))
	err = errors.New("Error 1062: Duplicate entry '1' for key 'PRIMARY'")
	assert.Equal(t, false, m.isSystemError(err))
}

func BenchmarkIsSystemError(b *testing.B) {
	m := &mysql{}
	err := errors.New("Error 1062: Duplicate entry '1' for key 'PRIMARY'")
	for i := 0; i < b.N; i++ {
		m.isSystemError(err)
	}
}

func BenchmarkSqlQuery(b *testing.B) {
	b.ReportAllocs()
	my := newMysql("hellofarm:halfquestfarm4321@tcp(127.0.0.1:3306)/ChatShard1?charset=utf8&timeout=4s", nil)
	my.Open()
	var err error
	for i := 0; i < b.N; i++ {
		_, e := my.Query("SELECT * FROM UserInfo WHERE uid=?", 1)
		if e != nil {
			err = e
		}
	}

	if err != nil {
		b.Log(err)
	}

}

func BenchmarkSqlExec(b *testing.B) {
	b.ReportAllocs()
	my := newMysql("hellofarm:halfquestfarm4321@tcp(127.0.0.1:3306)/ChatShard1?charset=utf8&timeout=4s", nil)
	my.Open()
	var err error
	for i := 0; i < b.N; i++ {
		_, _, e := my.ExecSql("UPDATE UserInfo SET gold=? WHERE uid=?", 12, 1)
		if e != nil {
			err = e
		}
	}

	if err != nil {
		b.Log(err)
	}

}
