package mysql

type ServerSelector interface {
	SetServers(servers ...string)
	PickServer(pool string, shardId int) (addr string, err error)
}
