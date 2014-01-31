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

type MemcacheClient struct {
	addr string

	conn   net.Conn
	buffer bufio.ReadWriter
}

type Result struct {
	Key   string
	Value []byte
	Flags uint16
	Cas   uint64
}

func newMemcacheClient() *MemcacheClient {
	return &MemcacheClient{}
}

func (this *MemcacheClient) Connect(addr string) (err error) {
	var network string
	if strings.Contains(addr, "/") {
		network = "unix"
	} else {
		network = "tcp"
	}

	var conn net.Conn
	conn, err = net.Dial(network, addr)
	if err != nil {
		return err
	}

	this.conn = conn
	this.buffer = bufio.ReadWriter{
		Reader: bufio.NewReader(this.conn),
		Writer: bufio.NewWriter(this.conn),
	}

	return nil
}

func (this *MemcacheClient) Close() {
	this.conn.Close()
	this.conn = nil
}

func (this *MemcacheClient) IsClosed() bool {
	return this.conn == nil
}

func (this *MemcacheClient) Get(keys ...string) (results []Result, err error) {
	defer handleError(&err)
	results = this.get("get", keys)
	return
}

func (this *MemcacheClient) Gets(keys ...string) (results []Result, err error) {
	defer handleError(&err)
	results = this.get("gets", keys)
	return
}

func (this *MemcacheClient) Set(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("set", key, flags, timeout, value, 0), nil
}

func (this *MemcacheClient) Add(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("add", key, flags, timeout, value, 0), nil
}

func (this *MemcacheClient) Replace(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("replace", key, flags, timeout, value, 0), nil
}

func (this *MemcacheClient) Append(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("append", key, flags, timeout, value, 0), nil
}

func (this *MemcacheClient) Prepend(key string, flags uint16, timeout uint64,
	value []byte) (stored bool, err error) {
	defer handleError(&err)
	return this.store("prepend", key, flags, timeout, value, 0), nil
}

func (this *MemcacheClient) Cas(key string, flags uint16, timeout uint64,
	value []byte, cas uint64) (stored bool, err error) {
	defer handleError(&err)
	return this.store("cas", key, flags, timeout, value, cas), nil
}

func (this *MemcacheClient) Delete(key string) (deleted bool, err error) {
	defer handleError(&err)
	// delete <key> [<time>] [noreply]\r\n
	this.writeStrings("delete", key, "\r\n")
	reply := this.readline()
	if strings.Contains(reply, "ERROR") {
		panic(newMemcacheError("Server error"))
	}

	return strings.HasPrefix(reply, "DELETED"), nil
}

func (this *MemcacheClient) get(cmd string, keys []string) (results []Result) {
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

func (this *MemcacheClient) store(cmd string, key string, flags uint16, timeout uint64,
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

func (this *MemcacheClient) writeString(s string) {
	if _, err := this.buffer.WriteString(s); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *MemcacheClient) writeStrings(strs ...string) {
	for _, s := range strs {
		this.writeString(s)
	}
}

func (this *MemcacheClient) write(b []byte) {
	if _, err := this.buffer.Write(b); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *MemcacheClient) flush() {
	if err := this.buffer.Flush(); err != nil {
		panic(newMemcacheError("%s", err))
	}
}

func (this *MemcacheClient) readline() string {
	this.flush()
	line, isPrefix, err := this.buffer.ReadLine()
	if isPrefix || err != nil {
		panic(newMemcacheError("Prefix: %v, %s", isPrefix, err))
	}

	return string(line)
}

func (this *MemcacheClient) read(count int) []byte {
	this.flush()
	buf := make([]byte, count)
	if _, err := io.ReadFull(this.buffer, buf); err != nil {
		panic(newMemcacheError("%s", err))
	}

	return buf
}
