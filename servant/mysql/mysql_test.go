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

// 242102251 ns/op	    3427 B/op	      40 allocs/op
// mysql_test.go:40: dial tcp 192.168.23.163:3306: i/o timeout
func BenchmarkSqlQuery(b *testing.B) {
	b.ReportAllocs()
	my := newMysql("hellofarm:halfquestfarm4321@tcp(192.168.23.163:3306)/UserShard1?charset=utf8&timeout=4s", 0, nil)
	my.Open()
	//my.db.SetMaxIdleConns(100)
	//my.db.SetMaxOpenConns(30)
	query := "SELECT * FROM UserInfo WHERE uid=?"
	for i := 0; i < b.N; i++ {
		rows, e := my.Query(query, 1)
		if e != nil {
			b.Log(e)
		} else {
			rows.Close()
		}

	}

}

// 1005303 ns/op	     275 B/op	      16 allocs/op
func Benc1hmarkSqlExec(b *testing.B) {
	b.ReportAllocs()
	my := newMysql("hellofarm:halfquestfarm4321@tcp(192.168.23.163:3306)/UserShard1?charset=utf8&timeout=4s", 0, nil)
	my.Open()
	//my.db.SetMaxIdleConns(100)
	//my.db.SetMaxOpenConns(1000)
	for i := 0; i < b.N; i++ {
		_, _, e := my.Exec("UPDATE UserInfo SET gold=? WHERE uid=?", 12, 1)
		if e != nil {
			b.Log(e)
		}
	}

}
