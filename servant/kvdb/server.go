package kvdb

import (
	"fmt"
	"os"
	"runtime"
)

type Server struct {
	servlets []*Servlet
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

	this.servlets = make([]*Servlet, 0)
	for i := 0; i < runtime.NumCPU(); i++ {
		servlet := NewServlet(fmt.Sprintf("%s/%d", this.path, i))
		servlet.Open()
		this.servlets = append(this.servlets, servlet)
	}

	return nil
}

func (this *Server) close() {
	if this.servlets != nil {
		for _, servlet := range this.servlets {
			servlet.Close()
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
