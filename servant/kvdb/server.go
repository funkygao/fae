package kvdb

import (
	"fmt"
	"os"
	"runtime"
)

type Server struct {
	shards   int
	servlets []*servlet
	path     string
}

func NewServer(path string, servletNum int) *Server {
	if servletNum == 0 {
		servletNum = runtime.NumCPU()
	}
	return &Server{path: path, shards: servletNum}
}

func (this *Server) Open() (err error) {
	this.close()

	err = this.createIfNotExists()
	if err != nil {
		return err
	}

	this.servlets = make([]*servlet, 0)
	for i := 0; i < this.shards; i++ {
		servlet := newServlet(fmt.Sprintf("%s/%d", this.path, i))
		if err := servlet.open(); err != nil {
			this.close()
			return err
		}

		this.servlets = append(this.servlets, servlet)
	}

	return nil
}

func (this *Server) Get(key []byte) (value []byte, err error) {
	servlet := this.servletByKey(key)
	servlet.lock()
	defer servlet.unlock()
	return servlet.get(key)
}

func (this *Server) Put(key []byte, value []byte) error {
	servlet := this.servletByKey(key)
	servlet.lock()
	defer servlet.unlock()
	return servlet.put(key, value)
}

func (this *Server) Delete(key []byte) error {
	servlet := this.servletByKey(key)
	servlet.lock()
	defer servlet.unlock()
	return servlet.delete(key)
}

func (this *Server) Count() (c int) {
	for _, s := range this.servlets {
		c += s.count()
	}
	return
}

func (this *Server) close() {
	if this.servlets != nil {
		for _, servlet := range this.servlets {
			servlet.close()
		}

		this.servlets = nil
	}
}

func (this *Server) createIfNotExists() (err error) {
	err = os.MkdirAll(this.path, DB_PERM)
	if err != nil {
		return err
	}

	return nil
}
