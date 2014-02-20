package kvdb

import (
	"hash/crc32"
)

func (this *Server) servletOwnerIndex(key []byte) uint32 {
	return crc32.ChecksumIEEE(key) % uint32(len(this.servlets))
}
