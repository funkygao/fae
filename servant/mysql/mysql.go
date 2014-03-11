package mysql

import (
	"database/sql"
	"fmt"
	log "github.com/funkygao/log4go"
	_ "github.com/go-sql-driver/mysql"
)

type mysql struct {
	dsn string
	db  *sql.DB
}

func newMysql(dsn string) *mysql {
	this := new(mysql)
	this.dsn = dsn

	// conn to db
	var err error
	this.db, err = sql.Open("mysql", this.dsn)
	this.checkError(err, dsn)

	return this
}

func (this mysql) String() string {
	return this.dsn
}

// sets the maximum number of connections in the idle connection pool
func (this *mysql) SetMaxIdleConns(n int) {
	this.db.SetMaxIdleConns(n)
}

func (this *mysql) checkError(err error, sql string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s, %s", this, err.Error(), sql))
	}
}

func (this *mysql) Query(query string, args ...interface{}) *sql.Rows {
	log.Debug("%s, args=%+v\n", query, args)

	rows, err := this.db.Query(query, args...)
	this.checkError(err, query)

	return rows
}

func (this *mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	log.Debug("%s, args=%+v\n", query, args)

	return this.db.QueryRow(query, args...)
}

func (this *mysql) ExecSql(query string, args ...interface{}) (afftectedRows int64) {
	log.Debug("%s, args=%+v\n", query, args)

	res, err := this.db.Exec(query, args...)
	this.checkError(err, query)

	afftectedRows, err = res.RowsAffected()
	this.checkError(err, query)

	return
}

func (this *mysql) Prepare(query string) *sql.Stmt {
	log.Debug(query)

	r, err := this.db.Prepare(query)
	this.checkError(err, query)
	return r
}

func (this *mysql) Close() error {
	return this.db.Close()
}

func (this *mysql) Db() *sql.DB {
	return this.db
}
