package kvdb

import (
	"fmt"
	"github.com/jmhodges/levigo"
	"os"
	"sync"
)

// A Servlet is a small wrapper around a single shard of a LevelDB data file.
type Servlet struct {
	path  string // data/0
	db    *levigo.DB
	mutex sync.Mutex
}

// NewServlet returns a new Servlet with a data shard stored at a given path.
func NewServlet(path string) *Servlet {
	return &Servlet{path: path}
}

// Opens the underlying LevelDB database and starts the message loop.
func (this *Servlet) Open() error {
	err := os.MkdirAll(this.path, DB_PERM)
	if err != nil {
		return err
	}

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(1 << 30)) // TODO config
	filter := levigo.NewBloomFilter(10)
	opts.SetFilterPolicy(filter)
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(this.path, opts)
	if err != nil {
		panic(fmt.Sprintf("Unable to open LevelDB database: %v", err))
	}
	this.db = db
	return nil
}

// Closes the underlying LevelDB database.
func (this *Servlet) Close() {
	if this.db != nil {
		this.db.Close()
	}
}

// Locks the entire servlet.
func (this *Servlet) Lock() {
	this.mutex.Lock()
}

// Unlocks the entire servlet.
func (this *Servlet) Unlock() {
	this.mutex.Unlock()
}

func (this *Servlet) Get(key []byte) (value []byte, err error) {
	ro := levigo.NewReadOptions()
	defer ro.Close()
	return this.db.Get(ro, key)
}

func (this *Servlet) Put(key []byte, value []byte) error {
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	return this.db.Put(wo, key, value)
}

func (this *Servlet) Delete(key []byte) error {
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	return this.db.Delete(wo, key)
}
