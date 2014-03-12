package mysql

type ServerSelector interface {
	PickServer(pool string, table string, shardId int) (*mysql, error)
}
