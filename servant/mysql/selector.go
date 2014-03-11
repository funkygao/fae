package mysql

type ServerSelector interface {
	PickServer(pool string, shardId int) (*mysql, error)
}
