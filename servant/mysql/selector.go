package mysql

type ServerSelector interface {
	ServerByBucket(bucket string) (*mysql, error)
	PickServer(pool string, table string, hintId int) (*mysql, error)
	Servers() []*mysql
}
