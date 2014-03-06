package servant

func (this *FunServantImpl) warmUp() {
	go this.mg.WarmUp()
	go this.mc.WarmUp()
}
