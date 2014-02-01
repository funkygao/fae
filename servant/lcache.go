package servant

import (
	log "code.google.com/p/log4go"
	"time"
)

func (this *FunServantImpl) LcSet(key string, value []byte) (r bool, err error) {
	this.t1 = time.Now()
	this.lc.Set(key, value)
	r = true

	log.Debug("lc_set key:%s value:%v %s", key, value, time.Since(this.t1))
	return
}

func (this *FunServantImpl) LcGet(key string) (r []byte, err error) {
	this.t1 = time.Now()
	result, ok := this.lc.Get(key)
	if !ok {
		err = errLcMissed
	} else {
		r = result.([]byte)
	}

	log.Debug("lc_get key:%s %s", key, time.Since(this.t1))
	return
}

func (this *FunServantImpl) LcDel(key string) (err error) {
	log.Debug("lc_delete key:%s", key)
	this.lc.Del(key)
	return
}
