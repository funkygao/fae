package mysql

import (
	"database/sql"
	log "github.com/funkygao/log4go"
	_ "github.com/go-sql-driver/mysql"
)

// A mysql conn to a single mysql instance
// Conn pool is natively supported by golang
type mysql struct {
	dsn string
	db  *sql.DB
}

func newMysql(dsn string) *mysql {
	this := new(mysql)
	this.dsn = dsn

	return this
}

func (this *mysql) Open() (err error) {
	this.db, err = sql.Open("mysql", this.dsn)
	return
}

func (this mysql) String() string {
	return this.dsn
}

// sets the maximum number of connections in the idle connection pool
func (this *mysql) SetMaxIdleConns(n int) {
	this.db.SetMaxIdleConns(n)
}

func (this *mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	log.Debug("%s, args=%+v\n", query, args)

	return this.db.Query(query, args...)
}

func (this *mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	log.Debug("%s, args=%+v\n", query, args)

	return this.db.QueryRow(query, args...)
}

func (this *mysql) ExecSql(query string, args ...interface{}) (afftectedRows int64, err error) {
	log.Debug("%s, args=%+v\n", query, args)

	var result sql.Result
	result, err = this.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	afftectedRows, err = result.RowsAffected()
	return

}

func (this *mysql) Close() error {
	return this.db.Close()
}
