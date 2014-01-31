package engine

import (
	log "code.google.com/p/log4go"
	"fmt"
	conf "github.com/daviddengcn/go-ljson-conf"
)

type configRpc struct {
	listenAddr string
	framed     bool
	protocol   string
}

func (this *configRpc) loadConfig(section *conf.Conf) {
	this.listenAddr = section.String("listen_addr", "")
	if this.listenAddr == "" {
		panic("Empty listen_addr")
	}
	this.framed = section.Bool("framed", false)
	this.protocol = section.String("protocol", "binary")

	log.Debug("rpc: %+v", *this)
}

type configMemcacheServer struct {
	host string
	port string
}

func (this *configMemcacheServer) loadConfig(section *conf.Conf) {
	this.host = section.String("host", "")
	if this.host == "" {
		panic("Empty memcache server host")
	}
	this.port = section.String("port", "")
	if this.port == "" {
		panic("Empty memcache server port")
	}

	log.Debug("memcache server: %+v", *this)
}

type configMemcache struct {
	hashStrategy string
	hashFunction string
	servers      map[string]*configMemcacheServer // key is host:port
}

func (this *configMemcache) loadConfig(cf *conf.Conf) {
	this.servers = make(map[string]*configMemcacheServer)
	this.hashStrategy = cf.String("hash_strategy", "standard")
	this.hashFunction = cf.String("hash_function", "crc32")
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(configMemcacheServer)
		server.loadConfig(section)
		this.servers[server.host+":"+server.port] = server
	}

	log.Debug("memcache: %+v", *this)
}

type configMongodbServer struct {
	shardName  string
	host       string
	port       int
	user, pass string
	db         string
	replicaSet string
}

func (this *configMongodbServer) loadConfig(section *conf.Conf) {
	this.shardName = section.String("shard_name", "")
	this.host = section.String("host", "")
	this.port = section.Int("port", 0)
	this.db = section.String("db", "")
	this.user = section.String("user", "")
	this.pass = section.String("pass", "")
	this.replicaSet = section.String("replicaSet", "")
	if this.host == "" ||
		this.port == 0 ||
		this.shardName == "" ||
		this.db == "" {
		panic("required filed")
	}

	log.Debug("mongodb server: %+v", *this)
}

type configMongodb struct {
	servers map[string]*configMongodbServer // key is shardName
}

func (this *configMongodb) loadConfig(cf *conf.Conf) {
	this.servers = make(map[string]*configMongodbServer)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(configMongodbServer)
		server.loadConfig(section)
		this.servers[server.shardName] = server
	}

	log.Debug("mongodb: %+v", *this)
}

type engineConfig struct {
	*conf.Conf

	httpListenAddr string

	rpc      *configRpc
	mongodb  *configMongodb
	memcache *configMemcache
}

func (this *Engine) LoadConfigFile() *Engine {
	log.Debug("Loading config file %s", this.configFile)

	config := new(engineConfig)
	var err error
	config.Conf, err = conf.Load(this.configFile)
	if err != nil {
		panic(err)
	}

	this.conf = config
	this.doLoadConfig()

	return this
}

func (this *Engine) doLoadConfig() {
	this.conf.httpListenAddr = this.conf.String("http_listen_addr", "")

	// rpc section
	this.conf.rpc = new(configRpc)
	section, err := this.conf.Section("rpc")
	if err != nil {
		panic(err)
	}
	this.conf.rpc.loadConfig(section)

	// mongodb section
	this.conf.mongodb = new(configMongodb)
	section, err = this.conf.Section("mongodb")
	if err != nil {
		panic(err)
	}
	this.conf.mongodb.loadConfig(section)

	// memcached section
	this.conf.memcache = new(configMemcache)
	section, err = this.conf.Section("memcache")
	if err != nil {
		panic(err)
	}
	this.conf.memcache.loadConfig(section)
}
