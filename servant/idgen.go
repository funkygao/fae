package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
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

func (this *FunServantImpl) milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context,
	flag int16) (r int64, backwards *rpc.TIdTimeBackwards, appErr error) {
	this.idgenMutex.Lock()
	defer this.idgenMutex.Unlock()

	var (
		did int64 = 0 // datacenter id
		wid int64 = 0 // worker id
	)

	ts := this.milliseconds()
	if ts < this.idLastTimestamp {
		backwards = rpc.NewTIdTimeBackwards()
		return
	}

	if this.idLastTimestamp == ts {
		this.idSeq = (this.idSeq + 1) & sequenceMask
		if this.idSeq == 0 {
			for ts <= this.idLastTimestamp {
				ts = this.milliseconds()
			}
		}
	} else {
		this.idSeq = 0
	}

	this.idLastTimestamp = ts

	r = ((ts - twepoch) << timestampLeftShift) |
		(did << datacenterIdShift) |
		(wid << workerIdShift) |
		this.idSeq

	return
}
