package mysql

type ServerSelector interface {
	PickServer(pool string, table string, hintId int) (*mysql, error)
	Servers() []*mysql
}
