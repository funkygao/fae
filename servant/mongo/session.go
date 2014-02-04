package mongo

import (
	"labix.org/v2/mgo"
)

type Session struct {
	*mgo.Session

	addr   string
	client *Client
}

func (this *Session) resumableError(err error) bool {
	switch err {
	case nil:
		return true
	}

	return false
}

func (this *Session) Recyle(err *error) {
	if err == nil || this.resumableError(*err) {
		this.client.putFreeConn(this.addr, this.Session)
	} else {
		this.Session.Close()
	}
}
