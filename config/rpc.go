package config

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ConfigRpc struct {
	ListenAddr             string
	SessionTimeout         time.Duration
	IoTimeout              time.Duration
	BufferSize             int // network IO read/write buffer
	Framed                 bool
	Protocol               string
	StatsOutputInterval    time.Duration
	PreforkMode            bool
	MaxOutstandingSessions int
	HostMaxCallPerMinute   int
}

func (this *ConfigRpc) LoadConfig(section *conf.Conf) {
	this.ListenAddr = section.String("listen_addr", "")
	if this.ListenAddr == "" {
		panic("Empty listen_addr")
	}

	this.SessionTimeout = section.Duration("session_timeout", 30*time.Second)
	this.IoTimeout = section.Duration("io_timeout", 2*time.Second)
	this.StatsOutputInterval = section.Duration("stats_output_interval", 10*time.Second)
	this.Framed = section.Bool("framed", false)
	this.BufferSize = section.Int("buffer_size", 4<<10)
	this.Protocol = section.String("protocol", "binary")
	this.PreforkMode = section.Bool("prefork_mode", false)
	this.MaxOutstandingSessions = section.Int("max_outstanding_sessions", 20000)
	this.HostMaxCallPerMinute = section.Int("host_max_call_per_minute", 100*60)

	log.Debug("rpc conf: %+v", *this)
}
