package config

import (
	"encoding/json"
	"fmt"
	"time"

	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type ConfigMysqlServer struct {
	Pool    string `json:"pool"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	DbName  string `json:"db"`
	Charset string `json:"charset"`

	conf *ConfigMysql
	dsn  string // cache of op result
}

func (this *ConfigMysqlServer) loadConfig(section *conf.Conf) {
	this.Pool = section.String("pool", "")
	this.Host = section.String("host", "")
	this.Port = section.String("port", "3306")
	this.DbName = section.String("db", "")
	this.User = section.String("username", "")
	this.Pass = section.String("password", "")
	this.Charset = section.String("charset", "utf8")
	if this.Host == "" ||
		this.Port == "" ||
		this.Pool == "" ||
		this.DbName == "" {
		panic("required field missing")
	}

	this.fillDsn()
}

func (this *ConfigMysqlServer) fillDsn() {
	this.dsn = ""
	if this.User != "" {
		this.dsn = this.User + ":"
		if this.Pass != "" {
			this.dsn += this.Pass
		}
	}
	this.dsn += fmt.Sprintf("@tcp(%s:%s)/%s?", this.Host, this.Port, this.DbName)
	if this.Charset != "" {
		this.dsn += "charset=" + this.Charset
	}
	if this.conf.Timeout > 0 {
		this.dsn += "&timeout=" + this.conf.Timeout.String()
	}
}

func (this *ConfigMysqlServer) DSN() string {
	return this.dsn
}

type ConfigMysql struct {
	ShardStrategy                string                        `json:"shard_stategy"`
	Timeout                      time.Duration                 `json:"timeout"`
	GlobalPools                  map[string]bool               `json:"global_pools"` // non-sharded pools
	MaxIdleTime                  time.Duration                 `json:"idle_timeout"`
	MaxIdleConnsPerServer        int                           `json:"max_idle_conns"`
	MaxConnsPerServer            int                           `json:"max_conns"`
	HeartbeatInterval            int                           `json:"-"`
	JsonMergeMaxOutstandingItems int                           `json:"-"`
	CachePrepareStmtMaxItems     int                           `json:"-"` // 0 means disabled
	AllowNullableColumns         bool                          `json:"-"`
	Breaker                      ConfigBreaker                 `json:"breaker"`
	Servers                      map[string]*ConfigMysqlServer `json:"pools"` // key is pool

	// cache related
	CacheStore            string `json:"cache_store"`
	CacheStoreRedisPool   string `json:"-"`
	CacheStoreMemMaxItems int    `json:"cache_cap"`
	CacheKeyHash          bool   `json:"cache_keyhash"`

	LookupCacheMaxItems int    `json:"lookup_cache_max_items"`
	LookupPool          string `json:"lookup_pool"`
	DefaultLookupTable  string `json:"default_lookup_table"`

	lookupTables conf.Conf
}

func (this *ConfigMysql) LoadConfig(cf *conf.Conf) {
	this.GlobalPools = make(map[string]bool)
	for _, p := range cf.StringList("global_pools", nil) {
		this.GlobalPools[p] = true
	}
	this.ShardStrategy = cf.String("shard_strategy", "standard")
	this.MaxIdleTime = cf.Duration("max_idle_time", 0)
	this.Timeout = cf.Duration("timeout", 10*time.Second)
	this.AllowNullableColumns = cf.Bool("allow_nullable_columns", true)
	this.MaxIdleConnsPerServer = cf.Int("max_idle_conns_per_server", 2)
	this.MaxConnsPerServer = cf.Int("max_conns_per_server",
		this.MaxIdleConnsPerServer*5)
	this.CachePrepareStmtMaxItems = cf.Int("cache_prepare_stmt_max_items", 0)
	this.HeartbeatInterval = cf.Int("heartbeat_interval", 120)
	this.CacheStore = cf.String("cache_store", "mem")
	this.CacheStoreMemMaxItems = cf.Int("cache_store_mem_max_items", 10<<20)
	this.CacheStoreRedisPool = cf.String("cache_store_redis_pool", "db_cache")
	this.CacheKeyHash = cf.Bool("cache_key_hash", false)
	this.DefaultLookupTable = cf.String("default_lookup_table", "")
	this.LookupPool = cf.String("lookup_pool", "ShardLookup")
	this.JsonMergeMaxOutstandingItems = cf.Int("json_merge_max_outstanding_items", 8<<20)
	this.LookupCacheMaxItems = cf.Int("lookup_cache_max_items", 1<<20)
	section, err := cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}
	section, err = cf.Section("lookup_tables")
	if err == nil {
		this.lookupTables = *section
	} else {
		panic(err)
	}

	this.Servers = make(map[string]*ConfigMysqlServer)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMysqlServer)
		server.conf = this
		server.loadConfig(section)
		this.Servers[server.Pool] = server
	}

	log.Debug("mysql conf: %+v", *this)
}

func (this *ConfigMysql) From(b []byte) error {
	err := json.Unmarshal(b, this)
	if err != nil {
		return err
	}

	// setup the internal structs
	for _, server := range this.Servers {
		server.conf = this
		server.fillDsn()
		if server.dsn == "" {
			return fmt.Errorf("empty mysql DSN for server: %s/%s", server.Pool, server.DbName)
		}
	}

	return nil
}

func (this *ConfigMysql) Enabled() bool {
	return len(this.Servers) > 0
}

func (this *ConfigMysql) Pools() (pools []string) {
	poolsMap := make(map[string]bool)
	for _, server := range this.Servers {
		poolsMap[server.Pool] = true
	}
	for poolName, _ := range poolsMap {
		pools = append(pools, poolName)
	}
	return
}

func (this *ConfigMysql) LookupTable(pool string) (lookupTable string) {
	return this.lookupTables.String(pool, this.DefaultLookupTable)
}
