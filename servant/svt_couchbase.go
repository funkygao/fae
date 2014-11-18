package servant

import (
	"github.com/couchbase/gomemcached"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// curl localhost:8091/pools/ | python -m json.tool
// curl localhost:8091/poolsStreaming/default?uuid=ee6009fb8f1ba20b3101a465455828ee

func (this *FunServantImpl) CbDel(ctx *rpc.Context, bucket string,
	key string) (appErr error) {
	const IDENT = "cb.del"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	appErr = b.Delete(key)
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx, "{b^%s k^%s",
		bucket, key)

	return
}

func (this *FunServantImpl) CbGet(ctx *rpc.Context, bucket string,
	key string) (r []byte, appErr error) {
	const IDENT = "cb.get"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	r, appErr = b.GetRaw(key) // FIXME 如果不存在，也会抛错，需要额外处理
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s} {r^%s}",
		bucket, key, string(r))

	return
}

// key can be up to 250 chars long, unique within a bucket
// val can be up to 25MB in size
func (this *FunServantImpl) CbSet(ctx *rpc.Context, bucket string,
	key string, val []byte, expire int32) (appErr error) {
	const IDENT = "cb.set"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	appErr = b.SetRaw(key, int(expire), val)
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s exp^%d}",
		bucket, key, string(val), expire)

	return
}

func (this *FunServantImpl) CbAdd(ctx *rpc.Context, bucket string,
	key string, val []byte, expire int32) (r bool, appErr error) {
	const IDENT = "cb.add"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	r, appErr = b.AddRaw(key, int(expire), val)
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s exp^%d} {r^%v}",
		bucket, key, string(val), expire, r)

	return
}

// append raw data to an existing item
func (this *FunServantImpl) CbAppend(ctx *rpc.Context, bucket string,
	key string, val []byte) (appErr error) {
	const IDENT = "cb.append"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	appErr = b.Append(key, val)
	if appErr != nil {
		log.Error("Q=%s %s: %s %s", IDENT, ctx.String(), key, appErr)
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%s v^%s}",
		bucket, key, string(val))

	return
}

// fetches multiple keys concurrently
func (this *FunServantImpl) CbGets(ctx *rpc.Context, bucket string,
	keys []string) (r map[string][]byte, appErr error) {
	const IDENT = "cb.gets"
	if this.cb == nil {
		appErr = ErrServantNotStarted
		return
	}

	profiler, err := this.getSession(ctx).startProfiler()
	if err != nil {
		appErr = err
		return
	}

	this.stats.inc(IDENT)

	b, err := this.cb.GetBucket(bucket)
	if err != nil {
		appErr = err
		return
	}

	var rv map[string]*gomemcached.MCResponse
	rv, appErr = b.GetBulk(keys)
	if appErr != nil {
		log.Error("Q=%s %s: %v %s", IDENT, ctx.String(), keys, appErr)
	} else {
		for k, data := range rv {
			r[k] = data.Body
		}
	}

	profiler.do(IDENT, ctx,
		"{b^%s k^%v v^%s} {r^%v}",
		bucket, keys, r)

	return
}
