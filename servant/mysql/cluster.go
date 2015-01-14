package mysql

import (
	"database/sql"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"time"
)

type MysqlCluster struct {
	selector ServerSelector
}

func New(cf *config.ConfigMysql) *MysqlCluster {
	this := new(MysqlCluster)
	switch cf.ShardStrategy {
	case "standard":
		this.selector = newStandardServerSelector(cf)

	case "vbucket":
		panic("vbucket mysql sharding not implemented")

	default:
		panic("unknown mysql sharding type: " + cf.ShardStrategy)
	}

	return this
}

func (this *MysqlCluster) Query(pool string, table string, hintId int,
	sql string, args []interface{}) (*sql.Rows, error) {
	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return nil, err
	}

	return my.Query(sql, args...)
}

func (this *MysqlCluster) Exec(pool string, table string, hintId int,
	sql string, args []interface{}) (afftectedRows int64,
	lastInsertId int64, err error) {
	this.selector.KickLookupCache(pool, hintId)

	my, err := this.selector.PickServer(pool, table, hintId)
	if err != nil {
		return 0, 0, err
	}

	return my.ExecSql(sql, args...)
}

func (this *MysqlCluster) KickLookupCache(pool string, hintId int) {
	this.selector.KickLookupCache(pool, hintId)
}

func (this *MysqlCluster) Warmup() {
	var (
		err error
		t1  = time.Now()
	)

	for _, m := range this.selector.Servers() {
		err = m.Ping()
		if err != nil {
			log.Error("Warmup mysql: %s", err)
			break
		}

	}

	if err != nil {
		log.Error("Mysql failed to warmup within %s: %s",
			time.Since(t1), err)
	} else {
		log.Debug("Mysql warmup within %s: %+v",
			time.Since(t1), this.selector)
	}
}
