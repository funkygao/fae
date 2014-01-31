/*
Memcache client
*/
package memcache

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	// 1MB
	MAX_VALUE_SIZE = 1000000
)

type Memcache struct {
	conn   net.Conn
	buffer bufio.ReadWriter
}

type Result struct {
	Key   string
	Value []byte
	Flags uint16
	Cas   uint64
}

func Connect(addr string) (mc *Memcache, err error) {
	var network string
	if strings.Contains(addr, "/") {
		network = "unix"
	} else {
		network = "tcp"
	}

	var nc net.Conn
	nc, err = net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return newMemcache(nc), nil
}

func newMemcache(nc net.Conn) *Memcache {
	return &Memcache{
		conn: nc,
		buffer: bufio.ReadWriter{
			Reader: bufio.NewReader(nc),
			Writer: bufio.NewWriter(nc),
		},
	}
}

func (this *Memcache) Close() {
	this.conn.Close()
	this.conn = nil
}

func (this *Memcache) IsClosed() bool {
	return this.conn == nil
}

func (this *Memcache) Get(keys ...string) (results []Result, err error) {
	defer handleError(&err)
	results = this.get("get", keys)
	return
}

func (this *Memcache) Gets(keys ...string) (results []Result, err error) {
	defer handleError(&err)
	results = this.get("gets", keys)
	return
}

func (this *Memcache) Set(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("set", key, flags, timeout, value, 0), nil
}

func (this *Memcache) Add(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("add", key, flags, timeout, value, 0), nil
}

func (this *Memcache) Replace(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("replace", key, flags, timeout, value, 0), nil
}

func (this *Memcache) Append(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("append", key, flags, timeout, value, 0), nil
}

func (this *Memcache) Prepend(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("prepend", key, flags, timeout, value, 0), nil
}

func (this *Memcache) Cas(key string, flags uint16, timeout uint64,
	value []byte, cas uint64) (stored bool, err error) {
	defer handleError(&err)
	return this.store("cas", key, flags, timeout, value, cas), nil
}

func (this *Memcache) Delete(key string) (deleted bool, err error) {
	defer handleError(&err)
	// delete <key> [<time>] [noreply]\r\n
	this.writeStrings("delete", key, "\r\n")
	reply := this.readline()
	if strings.Contains(reply, "ERROR") {
		panic(newMemcacheError("Server error"))
	}

	return strings.HasPrefix(reply, "DELETED"), nil
}

func (this *Memcache) get(cmd string, keys []string) (results []Result) {
	results = make([]Result, 0, len(keys))
	if len(keys) == 0 {
		return
	}

	// get(s) <key>\r\n
	this.writeString(cmd)
	for _, key := range keys {
		this.writeStrings(" ", key)
	}
	this.writeString("\r\n")

	header := this.readline()
	var result Result
	for strings.HasPrefix(header, "VALUE") {
		// VALUE <key> <flags> <bytes> [<cas unique>]\r\n
		chunks := strings.Split(header, " ")
		if len(chunks) < 4 {
			panic(newMemcacheError("Malformed response: %s", string(header)))
		}
		result.Key = chunks[1]
		flags64, err := strconv.ParseUint(chunks[2], 10, 16)
		if err != nil {
			panic(newMemcacheError("%v", err))
		}
		result.Flags = uint16(flags64)
		size, err := strconv.ParseUint(chunks[3], 10, 64)
		if err != nil {
			panic(newMemcacheError("%v", err))
		}
		if len(chunks) == 5 {
			result.Cas, err = strconv.ParseUint(chunks[4], 10, 64)
			if err != nil {
				panic(newMemcacheError("%v", err))
			}
		}

		// <data block>\r\n
		result.Value = this.read(int(size) + 2)[:size]
		results = append(results, result)

		header = this.readline()
	}

	if !strings.HasPrefix(header, "END") {
		panic(newMemcacheError("Malformed response: %s", string(header)))
	}

	return
}

func (this *Memcache) store(cmd string, key string, flags uint16, timeout uint64,
	value []byte, cas uint64) (stored bool) {
	if len(value) > MAX_VALUE_SIZE {
		return false
	}

	// <cmd> <key> <flags> <exptime> <bytes> [noreply]\r\n
	this.writeStrings(cmd, " ", key, " ")
	this.write(strconv.AppendUint(nil, uint64(flags), 10))
	this.writeString(" ")
	this.write(strconv.AppendUint(nil, timeout, 10))
	this.writeString(" ")
	this.write(strconv.AppendInt(nil, int64(len(value)), 10))
	if cas != 0 {
		this.writeString(" ")
		this.write(strconv.AppendUint(nil, cas, 10))
	}
	this.writeString("\r\n")

	// <data block>\r\n
	this.write(value)
	this.writeString("\r\n")

	reply := this.readline()
	if strings.Contains(reply, "ERROR") {
		panic(newMemcacheError("Server error"))
	}

	return strings.HasPrefix(reply, "STORED")
}

func (this *Memcache) writeString(s string) {
	if _, err := this.buffer.WriteString(s); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *Memcache) writeStrings(strs ...string) {
	for _, s := range strs {
		this.writeString(s)
	}
}

func (this *Memcache) write(b []byte) {
	if _, err := this.buffer.Write(b); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *Memcache) flush() {
	if err := this.buffer.Flush(); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *Memcache) readline() string {
	this.flush()
	line, isPrefix, err := this.buffer.ReadLine()
	if isPrefix || err != nil {
		panic(newMemcacheError("Prefix: %v, %s", isPrefix, err))
	}

	return string(line)
}

func (this *Memcache) read(count int) []byte {
	this.flush()
	buf := make([]byte, count)
	if _, err := io.ReadFull(this.buffer, buf); err != nil {
		panic(newMemcacheError("%s", err))
	}

	return buf
}
