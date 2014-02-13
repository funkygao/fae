package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"sync"
	"time"
)

const (
	workerIdBits       = uint64(5)
	datacenterIdBits   = uint64(5)
	maxWorkerId        = int64(-1) ^ (int64(-1) << workerIdBits)
	maxDatacenterId    = int64(-1) ^ (int64(-1) << datacenterIdBits)
	sequenceBits       = uint64(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)

	// Tue, 21 Mar 2006 20:50:14.000 GMT
	twepoch = int64(1288834974657)
)

type IdGenerator struct {
	mutex         sync.Mutex
	seq           int64
	lastTimestamp int64
}

func NewIdGenerator() (this *IdGenerator) {
	this = new(IdGenerator)
	return
}

func (this *IdGenerator) milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (this *IdGenerator) Next() (int64, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var (
		did int64 = 0 // datacenter id
		wid int64 = 0 // worker id
	)

	ts := this.milliseconds()
	if ts < this.lastTimestamp {
		return 0, rpc.NewTIdTimeBackwards()
	}

	if this.lastTimestamp == ts {
		this.seq = (this.seq + 1) & sequenceMask
		if this.seq == 0 {
			for ts <= this.lastTimestamp {
				ts = this.milliseconds()
			}
		}
	} else {
		this.seq = 0
	}

	this.lastTimestamp = ts

	r := ((ts - twepoch) << timestampLeftShift) |
		(did << datacenterIdShift) |
		(wid << workerIdShift) |
		this.seq
	return r, nil
}

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context,
	flag int16) (r int64, backwards *rpc.TIdTimeBackwards, appErr error) {
	r, appErr = this.idgen.Next()
	if appErr != nil {
		backwards = appErr.(*rpc.TIdTimeBackwards)
		appErr = nil
	}
	return
}
