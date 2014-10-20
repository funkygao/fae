package servant

func (this *FunServantImpl) warmUp() {
	if this.mg != nil {
		go this.mg.Warmup()
	}

	if this.mc != nil {
		go this.mc.Warmup()
	}

	// TODO warmup mysql

}
