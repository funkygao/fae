package config

import (
	"github.com/funkygao/golib/ip"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"strings"
	"time"
)

type ConfigProxy struct {
	PoolCapacity       int
	IdleTimeout        time.Duration
	IoTimeout          time.Duration
	BorrowTimeout      time.Duration
	DiagnosticInterval time.Duration
	SelfAddr           string
	TcpNoDelay         bool
	BufferSize         int
}

func NewDefaultProxy() *ConfigProxy {
	return &ConfigProxy{
		PoolCapacity:       10,
		IdleTimeout:        0,
		SelfAddr:           ":0",
		IoTimeout:          time.Second * 10,
		BorrowTimeout:      time.Second * 10,
		DiagnosticInterval: time.Second * 5,
		TcpNoDelay:         true,
		BufferSize:         4 << 10,
	}
}

func (this *ConfigProxy) LoadConfig(selfAddr string, cf *conf.Conf) {
	if selfAddr == "" {
		panic("proxy self addr unknown")
	}
	this.PoolCapacity = cf.Int("pool_capacity", 10)
	this.IdleTimeout = cf.Duration("idle_timeout", 0)
	this.IoTimeout = cf.Duration("io_timeout", time.Second*10)
	this.BorrowTimeout = cf.Duration("borrow_timeout", time.Second*10)
	this.DiagnosticInterval = cf.Duration("diagnostic_interval", time.Second*5)
	this.TcpNoDelay = cf.Bool("tcp_nodelay", true)
	this.BufferSize = cf.Int("buffer_size", 4<<10)
	this.SelfAddr = selfAddr
	parts := strings.SplitN(this.SelfAddr, ":", 2)
	if parts[0] == "" {
		// auto get local ip when self_addr like ":9001"
		this.SelfAddr = ip.LocalIpv4Addrs()[0] + ":" + parts[1]
	}

	log.Debug("proxy conf: %+v", *this)
}

func (this *ConfigProxy) Enabled() bool {
	return this.SelfAddr != ""
}
