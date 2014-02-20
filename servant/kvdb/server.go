package kvdb

import (
	"fmt"
	"os"
	"runtime"
)

type Server struct {
	servlets []*servlet
	path     string
}

func NewServer(path string) *Server {
	return &Server{path: path}
}

func (this *Server) Open() (err error) {
	this.close()

	err = this.createIfNotExists()
	if err != nil {
		return err
	}

	this.servlets = make([]*servlet, 0)
	for i := 0; i < runtime.NumCPU(); i++ {
		servlet := newServlet(fmt.Sprintf("%s/%d", this.path, i))
		servlet.open()
		this.servlets = append(this.servlets, servlet)
	}

	return nil
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
