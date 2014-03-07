package servant

func (this *FunServantImpl) warmUp() {
	go this.mg.Warmup()
	go this.mc.Warmup()
}
