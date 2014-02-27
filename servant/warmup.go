package servant

// FIXME
func (this *FunServantImpl) warmUp() {
	this.mg.WarmUp()
	this.mc.WarmUp()
}
