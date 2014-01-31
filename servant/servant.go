package servant

type FunServantImpl struct {
}

func NewFunServant() (this *FunServantImpl) {
	this = new(FunServantImpl)
	return
}

func init() {

}
