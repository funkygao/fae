package servant

import (
	"github.com/couchbase/gomemcached"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// curl localhost:8091/pools/ | python -m json.tool
// curl localhost:8091/poolsStreaming/default?uuid=ee6009fb8f1ba20b3101a465455828ee

func (this *FunServantImpl) CbDel(ctx *rpc.Context, bucket string,
	key string) (r bool, ex error) {
	const IDENT = "cb.del"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	ex = b.Delete(key)
	if ex != nil {
		r = false

		if e, ok := ex.(*gomemcached.MCResponse); ok && e.Status == gomemcached.KEY_ENOENT {
			ex = nil
		} else {
			// unexpected err
			log.Error("Q=%s %s %s: %s", IDENT, ctx.String(), key, ex.Error())
		}
	} else {
		// found this item, and deleted successfully
		r = true
	}

	profiler.do(IDENT, ctx, "{b^%s k^%s} {r^%v}",
		bucket, key, r)

	return
}

func (this *FunServantImpl) CbGet(ctx *rpc.Context, bucket string,
	key string) (r *rpc.TCouchbaseData, ex error) {
	const IDENT = "cb.get"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	r = rpc.NewTCouchbaseData()
	var data []byte
	data, ex = b.GetRaw(key)
	if ex != nil {
		r.Missed = true

		if e, ok := ex.(*gomemcached.MCResponse); ok && e.Status == gomemcached.KEY_ENOENT {
			ex = nil
		} else {
			log.Error("Q=%s %s %s: %s", IDENT, ctx.String(), key, ex.Error())
		}
	} else {
		r.Data = data
		r.Missed = false
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s} {r^%s}",
		bucket, key, string(r.Data))

	return
}

// key can be up to 250 chars long, unique within a bucket
// val can be up to 25MB in size
func (this *FunServantImpl) CbSet(ctx *rpc.Context, bucket string,
	key string, val []byte, expire int32) (ex error) {
	const IDENT = "cb.set"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	ex = b.SetRaw(key, int(expire), val)
	if ex != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, ex)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s exp^%d}",
		bucket, key, string(val), expire)

	return
}

func (this *FunServantImpl) CbAdd(ctx *rpc.Context, bucket string,
	key string, val []byte, expire int32) (r bool, ex error) {
	const IDENT = "cb.add"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	r, ex = b.AddRaw(key, int(expire), val)
	if ex != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, ex)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s exp^%d} {r^%v}",
		bucket, key, string(val), expire, r)

	return
}

// append raw data to an existing item
func (this *FunServantImpl) CbAppend(ctx *rpc.Context, bucket string,
	key string, val []byte) (ex error) {
	const IDENT = "cb.append"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	ex = b.Append(key, val)
	if ex != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, ex)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s}",
		bucket, key, string(val))

	return
}

// fetches multiple keys concurrently
func (this *FunServantImpl) CbGets(ctx *rpc.Context, bucket string,
	keys []string) (r map[string][]byte, ex error) {
	const IDENT = "cb.gets"
	if this.cb == nil {
		ex = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		ex = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		ex = err
		return
	}

	var rv map[string]*gomemcached.MCResponse
	rv, ex = b.GetBulk(keys)
	r = make(map[string][]byte)
	if ex != nil {
		log.Error("Q=%s %s: %v %s", IDENT, ctx.String(), keys, ex)
	} else {
		for k, data := range rv {
			r[k] = data.Body
		}
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%v} {r^%d}",
		bucket, keys, len(r))

	return
}
