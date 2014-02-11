package mongo

import (
	"github.com/funkygao/fae/config"
	"labix.org/v2/mgo"
)

type Session struct {
	*mgo.Session

	server *config.ConfigMongodbServer
	client *Client
}

func (this *Session) DB() *mgo.Database {
	return this.Session.DB(this.server.DbName)
}

func (this *Session) DbName() string {
	return this.server.DbName
}

func (this *Session) Recyle(err *error) {
	if err == nil || this.resumableError(*err) {
		// reusable session(connection)
		this.client.putFreeConn(this.server.Url(), this.Session)
	} else {
		// kill this session
		this.Session.Close()
	}
}

func (this *Session) resumableError(err error) bool {
	switch err {
	case nil, mgo.ErrNotFound:
		return true
	}

	return false
}
