package mysql

import (
	"database/sql"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
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
	if bc == nil {
		bc = &config.ConfigBreaker{
			FailureAllowance: 5,
			RetryTimeout:     time.Second * 10,
		}
	}
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

func (this *mysql) String() string {
	return this.dsn
}

func (this *mysql) Query(query string, args ...interface{}) (rows *sql.Rows,
	err error) {
	if this.db == nil {
		return nil, ErrNotOpen
	}
	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	rows, err = this.db.Query(query, args...)
	if err != nil {
		// func (me *MySQLError) Error() string {
		//     fmt.Sprintf("Error %d: %s", me.Number, me.Message)
		// }
		// Error 1054: Unknown column 'curve_internal_id' in 'field list'
		if this.isSystemError(err) {
			this.breaker.Fail()
		}
	} else {
		this.breaker.Succeed()
	}

	return
}

func (this *mysql) ExecSql(query string, args ...interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	if this.db == nil {
		return 0, 0, ErrNotOpen
	}
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

func (this *mysql) isSystemError(err error) bool {
	// http://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html
	// mysql error code is always 4 digits
	const (
		// Error 1054: Unknown column 'curve_internal_id' in 'field list'
		mysqlErrnoUnknownColumn = "1054"

		// Error 1062: Duplicate entry '1' for key 'PRIMARY'
		mysqlErrnoDupEntry = "1062"
	)

	// "Error %d:" skip the leading 6 chars: "Error "
	var errcode = err.Error()[6:] // TODO confirm mysql err always "Error %d: %s"
	switch {
	case strings.HasPrefix(errcode, mysqlErrnoUnknownColumn):
		return false

	case strings.HasPrefix(errcode, mysqlErrnoDupEntry):
		return false

	default:
		return true
	}
}
