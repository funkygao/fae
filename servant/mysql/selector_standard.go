package mysql

type StandardServerSelector struct {
}

func (this *StandardServerSelector) SetServers(servers ...string) {

}

func (this *StandardServerSelector) PickServer(pool string,
	shardId int) (addr string, err error) {
	return
}
