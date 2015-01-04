package mysql

type ServerSelector interface {
	ServerByBucket(bucket string) (*mysql, error)
	PickServer(pool string, table string, hintId int) (*mysql, error)
	KickLookupCache(pool string, hintId int)
	Servers() []*mysql
}
