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
	DiagnosticInterval time.Duration
	IoTimeout          time.Duration
	BorrowMaxSeconds   int
	SelfAddr           string
	TcpNoDelay         bool
	BufferSize         int

	enabled bool
}

func (this *ConfigProxy) LoadConfig(cf *conf.Conf) {
	this.PoolCapacity = cf.Int("pool_capacity", 10)
	this.IdleTimeout = cf.Duration("idle_timeout", 600*time.Second)
	this.SelfAddr = cf.String("self_addr", "")
	this.IoTimeout = cf.Duration("io_timeout", time.Second*10)
	this.BorrowMaxSeconds = cf.Int("borrow_max_seconds", 10)
	this.DiagnosticInterval = cf.Duration("diagnostic_interval", time.Second*30)
	this.TcpNoDelay = cf.Bool("tcp_nodelay", true)
	this.BufferSize = cf.Int("buffer_size", 4<<10)
	if this.SelfAddr == "" {
		log.Warn("proxy conf: %+v, empty self_addr: proxy disabled", *this)
		this.enabled = false
	} else {
		parts := strings.SplitN(this.SelfAddr, ":", 2)
		if parts[0] == "" {
			// auto get local ip when self_addr like ":9001"
			this.SelfAddr = ip.LocalIpv4Addrs()[0] + ":" + parts[1]
		}

		this.enabled = true
		log.Debug("proxy conf: %+v", *this)
	}
}

func (this *ConfigProxy) Enabled() bool {
	return this.enabled
}

func (this *ConfigProxy) Disable() {
	this.enabled = false
}
