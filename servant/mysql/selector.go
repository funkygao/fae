package mysql

type ServerSelector interface {
	ServerByBucket(bucket string) (*mysql, error)
	PickServer(pool string, table string, hintId int) (*mysql, error)
	Servers() []*mysql
	PoolServers(pool string) []*mysql
	KickLookupCache(pool string, hintId int)
}
