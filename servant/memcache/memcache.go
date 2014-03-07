package memcache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// mc client for a group of memcached instances
type Client struct {
	conf *config.ConfigMemcache

	selector ServerSelector

	lk            sync.Mutex
	breakers      map[net.Addr]*breaker.Consecutive
	freeconns     map[net.Addr][]*conn
	throttleConns map[net.Addr]chan interface{}
}

func newClient(cf *config.ConfigMemcache) (this *Client) {
	this = new(Client)
	this.conf = cf
	this.breakers = make(map[net.Addr]*breaker.Consecutive)
	this.throttleConns = make(map[net.Addr]chan interface{})

	switch cf.HashStrategy {
	case ConstistentHashStrategy:
		this.selector = new(ConsistentServerSelector)

	default:
		this.selector = new(StandardServerSelector)
	}

	if err := this.selector.SetServers(cf.ServerList()...); err != nil {
		panic(err)
	}

	return
}

func (this *Client) Warmup() {
	var (
		cn  *conn
		err error
		t1  = time.Now()
	)
	for retries := 0; retries < 3; retries++ {
		for _, addr := range this.selector.ServerList() {
			cn, err = this.getConn(addr)
			if err != nil {
				log.Error("Warmup %v fail: %s", addr, err)
				break
			} else {
				cn.condRelease(&err)
			}
		}

		if err == nil {
			// ok, needn't retry
			break
		}
	}

	if err == nil {
		log.Trace("Memcache warmed up within %s: %+v",
			time.Since(t1), this.freeconns)
	} else {
		log.Error("Memcache failed to warm up within %s: %s",
			time.Since(t1), err)
	}
}

func (this *Client) FreeConnMap() map[string][]*conn {
	this.lk.Lock()
	defer this.lk.Unlock()
	ret := make(map[string][]*conn)
	for addr, val := range this.freeconns {
		ret[addr.String()] = val
	}
	return ret
}

func (this *Client) putFreeConn(addr net.Addr, cn *conn) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if this.freeconns == nil {
		this.freeconns = make(map[net.Addr][]*conn)
	}
	freelist := this.freeconns[addr]
	if len(freelist) >= this.conf.MaxIdleConnsPerServer {
		cn.nc.Close()
		return
	}
	this.freeconns[addr] = append(freelist, cn)
}

func (this *Client) getFreeConn(addr net.Addr) (cn *conn, ok bool) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if _, present := this.breakers[addr]; !present {
		this.breakers[addr] = &breaker.Consecutive{
			FailureAllowance: this.conf.Breaker.FailureAllowance,
			RetryTimeout:     this.conf.Breaker.RetryTimeout}
		this.throttleConns[addr] = make(chan interface{}, this.conf.MaxConnsPerServer)
	}

	if this.freeconns == nil {
		return nil, false
	}
	freelist, ok := this.freeconns[addr]
	if !ok || len(freelist) == 0 {
		return nil, false
	}
	cn = freelist[len(freelist)-1]
	this.freeconns[addr] = freelist[:len(freelist)-1]
	return cn, true
}

func (this *Client) dial(addr net.Addr) (net.Conn, error) {
	if this.breakers[addr].Open() {
		return nil, ErrCircuitOpen
	}

	this.throttleConns[addr] <- true
	defer func() {
		// release throttle
		<-this.throttleConns[addr]
	}()

	type connError struct {
		cn  net.Conn
		err error
	}
	ch := make(chan connError)
	go func() {
		nc, err := net.Dial(addr.Network(), addr.String())
		ch <- connError{nc, err}
	}()
	select {
	case ce := <-ch:
		if ce.err != nil {
			this.breakers[addr].Fail()
		} else {
			this.breakers[addr].Succeed()
		}
		return ce.cn, ce.err
	case <-time.After(this.conf.Timeout):
		// Too slow. Fall through.
	}
	// Close the conn if it does end up finally coming in
	go func() {
		ce := <-ch
		if ce.err == nil {
			ce.cn.Close()
		}
	}()
	this.breakers[addr].Fail()
	return nil, &ConnectTimeoutError{addr}
}

func (this *Client) getConn(addr net.Addr) (*conn, error) {
	cn, ok := this.getFreeConn(addr)
	if ok {
		cn.extendDeadline()
		return cn, nil
	}
	nc, err := this.dial(addr)
	if err != nil {
		return nil, err
	}

	cn = &conn{
		nc:     nc,
		addr:   addr,
		rw:     bufio.NewReadWriter(bufio.NewReader(nc), bufio.NewWriter(nc)),
		client: this,
	}
	cn.extendDeadline()
	return cn, nil
}

func (this *Client) onItem(item *Item, fn func(*Client, *bufio.ReadWriter, *Item) error) error {
	addr, err := this.selector.PickServer(item.Key)
	if err != nil {
		return err
	}
	cn, err := this.getConn(addr)
	if err != nil {
		return err
	}
	defer cn.condRelease(&err)
	if err = fn(this, cn.rw, item); err != nil {
		return err
	}
	return nil
}

// Get gets the item for the given key. ErrCacheMiss is returned for a
// memcache cache miss. The key must be at most 250 bytes in length.
func (this *Client) Get(key string) (item *Item, err error) {
	err = this.withKeyAddr(key, func(addr net.Addr) error {
		return this.getFromAddr(addr, []string{key}, func(it *Item) { item = it })
	})
	if err == nil && item == nil {
		err = ErrCacheMiss
	}
	return
}

func (this *Client) withKeyAddr(key string, fn func(net.Addr) error) (err error) {
	if !legalKey(key) {
		return ErrMalformedKey
	}
	addr, err := this.selector.PickServer(key)
	if err != nil {
		return err
	}
	return fn(addr)
}

func (this *Client) withAddrRw(addr net.Addr, fn func(*bufio.ReadWriter) error) (err error) {
	cn, err := this.getConn(addr)
	if err != nil {
		return err
	}
	defer cn.condRelease(&err)
	return fn(cn.rw)
}

func (this *Client) withKeyRw(key string, fn func(*bufio.ReadWriter) error) error {
	return this.withKeyAddr(key, func(addr net.Addr) error {
		return this.withAddrRw(addr, fn)
	})
}

func (this *Client) getFromAddr(addr net.Addr, keys []string, cb func(*Item)) error {
	return this.withAddrRw(addr, func(rw *bufio.ReadWriter) error {
		if _, err := fmt.Fprintf(rw, "gets %s\r\n", strings.Join(keys, " ")); err != nil {
			return err
		}
		if err := rw.Flush(); err != nil {
			return err
		}
		if err := parseGetResponse(rw.Reader, cb); err != nil {
			return err
		}
		return nil
	})
}

// GetMulti is a batch version of Get. The returned map from keys to
// items may have fewer elements than the input slice, due to memcache
// cache misses. Each key must be at most 250 bytes in length.
// If no error is returned, the returned map will also be non-nil.
func (this *Client) GetMulti(keys []string) (map[string]*Item, error) {
	var lk sync.Mutex
	m := make(map[string]*Item)
	addItemToMap := func(it *Item) {
		lk.Lock()
		defer lk.Unlock()
		m[it.Key] = it
	}

	keyMap := make(map[net.Addr][]string)
	for _, key := range keys {
		if !legalKey(key) {
			return nil, ErrMalformedKey
		}
		addr, err := this.selector.PickServer(key)
		if err != nil {
			return nil, err
		}
		keyMap[addr] = append(keyMap[addr], key)
	}

	ch := make(chan error, buffered)
	for addr, keys := range keyMap {
		go func(addr net.Addr, keys []string) {
			ch <- this.getFromAddr(addr, keys, addItemToMap)
		}(addr, keys)
	}

	var err error
	for _ = range keyMap {
		if ge := <-ch; ge != nil {
			err = ge
		}
	}
	return m, err
}

// Set writes the given item, unconditionally.
func (this *Client) Set(item *Item) error {
	return this.onItem(item, (*Client).set)
}

func (this *Client) set(rw *bufio.ReadWriter, item *Item) error {
	return this.populateOne(rw, "set", item)
}

// Add writes the given item, if no value already exists for its
// key. ErrNotStored is returned if that condition is not met.
func (this *Client) Add(item *Item) error {
	return this.onItem(item, (*Client).add)
}

func (this *Client) add(rw *bufio.ReadWriter, item *Item) error {
	return this.populateOne(rw, "add", item)
}

// CompareAndSwap writes the given item that was previously returned
// by Get, if the value was neither modified or evicted between the
// Get and the CompareAndSwap calls. The item's Key should not change
// between calls but all other item fields may differ. ErrCASConflict
// is returned if the value was modified in between the
// calls. ErrNotStored is returned if the value was evicted in between
// the calls.
func (this *Client) CompareAndSwap(item *Item) error {
	return this.onItem(item, (*Client).cas)
}

func (this *Client) cas(rw *bufio.ReadWriter, item *Item) error {
	return this.populateOne(rw, "cas", item)
}

func (this *Client) populateOne(rw *bufio.ReadWriter, verb string, item *Item) error {
	if !legalKey(item.Key) {
		return ErrMalformedKey
	}
	var err error
	if verb == "cas" {
		_, err = fmt.Fprintf(rw, "%s %s %d %d %d %d\r\n",
			verb, item.Key, item.Flags, item.Expiration, len(item.Value), item.casid)
	} else {
		_, err = fmt.Fprintf(rw, "%s %s %d %d %d\r\n",
			verb, item.Key, item.Flags, item.Expiration, len(item.Value))
	}
	if err != nil {
		return err
	}
	if _, err = rw.Write(item.Value); err != nil {
		return err
	}
	if _, err := rw.Write(crlf); err != nil {
		return err
	}
	if err := rw.Flush(); err != nil {
		return err
	}
	line, err := rw.ReadSlice('\n')
	if err != nil {
		return err
	}
	switch {
	case bytes.Equal(line, resultStored):
		return nil
	case bytes.Equal(line, resultNotStored):
		return ErrNotStored
	case bytes.Equal(line, resultExists):
		return ErrCASConflict
	case bytes.Equal(line, resultNotFound):
		return ErrCacheMiss
	}
	return fmt.Errorf("memcache: unexpected response line from %q: %q", verb, string(line))
}

// Delete deletes the item with the provided key. The error ErrCacheMiss is
// returned if the item didn't already exist in the cache.
func (this *Client) Delete(key string) error {
	return this.withKeyRw(key, func(rw *bufio.ReadWriter) error {
		return writeExpectf(rw, resultDeleted, "delete %s\r\n", key)
	})
}

// Increment atomically increments key by delta. The return value is
// the new value after being incremented or an error. If the value
// didn't exist in memcached the error is ErrCacheMiss. The value in
// memcached must be an decimal number, or an error will be returned.
// On 64-bit overflow, the new value wraps around.
func (this *Client) Increment(key string, delta uint64) (newValue uint64, err error) {
	return this.incrDecr("incr", key, delta)
}

// Decrement atomically decrements key by delta. The return value is
// the new value after being decremented or an error. If the value
// didn't exist in memcached the error is ErrCacheMiss. The value in
// memcached must be an decimal number, or an error will be returned.
// On underflow, the new value is capped at zero and does not wrap
// around.
func (this *Client) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return this.incrDecr("decr", key, delta)
}

func (this *Client) incrDecr(verb, key string, delta uint64) (uint64, error) {
	var val uint64
	err := this.withKeyRw(key, func(rw *bufio.ReadWriter) error {
		line, err := writeReadLine(rw, "%s %s %d\r\n", verb, key, delta)
		if err != nil {
			return err
		}
		switch {
		case bytes.Equal(line, resultNotFound):
			return ErrCacheMiss
		case bytes.HasPrefix(line, resultClientErrorPrefix):
			errMsg := line[len(resultClientErrorPrefix) : len(line)-2]
			return errors.New("memcache: client error: " + string(errMsg))
		}
		val, err = strconv.ParseUint(string(line[:len(line)-2]), 10, 64)
		if err != nil {
			return err
		}
		return nil
	})

	return val, err
}
