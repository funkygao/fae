package servant

func (this *FunServantImpl) warmUp() {
	this.mg.WarmUp()
	this.mc.WarmUp()
}
