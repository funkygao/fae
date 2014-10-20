package mysql

import (
	"database/sql"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	_ "github.com/go-sql-driver/mysql"
)

// A mysql conn to a single mysql instance
// Conn pool is natively supported by golang
type mysql struct {
	dsn     string
	db      *sql.DB
	breaker *breaker.Consecutive
}

func newMysql(dsn string, bc *config.ConfigBreaker) *mysql {
	this := new(mysql)
	this.dsn = dsn
	this.breaker = &breaker.Consecutive{
		FailureAllowance: bc.FailureAllowance,
		RetryTimeout:     bc.RetryTimeout}

	return this
}

func (this *mysql) Open() (err error) {
	this.db, err = sql.Open("mysql", this.dsn)
	return
}

func (this *mysql) Ping() error {
	if this.db == nil {
		return ErrNotOpen
	}

	return this.db.Ping()
}

func (this mysql) String() string {
	return this.dsn
}

func (this *mysql) Query(query string, args ...interface{}) (rows *sql.Rows,
	err error) {
	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	rows, err = this.db.Query(query, args...)
	if err != nil {
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

	return
}

func (this *mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	return this.db.QueryRow(query, args...)
}

func (this *mysql) ExecSql(query string, args ...interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	if this.breaker.Open() {
		return 0, 0, ErrCircuitOpen
	}

	var result sql.Result
	result, err = this.db.Exec(query, args...)
	if err != nil {
		this.breaker.Fail()
		return 0, 0, err
	}

	afftectedRows, err = result.RowsAffected()
	if err != nil {
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

	lastInsertId, _ = result.LastInsertId()
	return

}
