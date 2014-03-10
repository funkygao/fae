package kvdb

import (
	"github.com/jmhodges/levigo"
	"os"
	"sync"
)

// A servlet is a small wrapper around a single shard of a LevelDB data file.
type servlet struct {
	path  string // data/0
	db    *levigo.DB
	mutex sync.Mutex
}

// Newservlet returns a new servlet with a data shard stored at a given path.
func newServlet(path string) *servlet {
	return &servlet{path: path}
}

// Opens the underlying LevelDB database and starts the message loop.
func (this *servlet) open() error {
	err := os.MkdirAll(this.path, DB_PERM)
	if err != nil {
		return err
	}

	opts := levigo.NewOptions()
	defer opts.Close()
	opts.SetCache(levigo.NewLRUCache(1 << 30)) // TODO config
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)
	opts.SetCreateIfMissing(true)
	opts.SetWriteBufferSize(1 << 20)
	opts.SetBlockSize(1024)
	opts.SetBlockRestartInterval(8)
	opts.SetCompression(levigo.SnappyCompression)
	db, err := levigo.Open(this.path, opts)
	if err != nil {
		return err
	}

	this.db = db
	return nil
}

// Closes the underlying LevelDB database.
func (this *servlet) close() {
	if this.db != nil {
		this.db.Close()
	}
}

// Locks the entire servlet.
func (this *servlet) lock() {
	this.mutex.Lock()
}

// Unlocks the entire servlet.
func (this *servlet) unlock() {
	this.mutex.Unlock()
}

func (this *servlet) get(key []byte) (value []byte, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	ro.SetVerifyChecksums(true)
	ro.SetFillCache(true)
	return this.db.Get(ro, key)
}

func (this *servlet) put(key []byte, value []byte) error {
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	wo.SetSync(false)
	return this.db.Put(wo, key, value)
}

func (this *servlet) delete(key []byte) error {
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	return this.db.Delete(wo, key)
}

func (this *servlet) count() (c int) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	it := this.db.NewIterator(ro)
	defer it.Close()
	it.SeekToFirst()
	for it = it; it.Valid(); it.Next() {
		c += 1
	}
	return
}
