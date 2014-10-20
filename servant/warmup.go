package servant

func (this *FunServantImpl) warmUp() {
	if this.mg != nil {
		go this.mg.Warmup()
	}

	if this.mc != nil {
		go this.mc.Warmup()
	}

	if this.my != nil {
		this.my.Warmup()
	}

}
