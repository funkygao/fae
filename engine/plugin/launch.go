package plugin

func Start() {
	for _, r := range runners {
		go r.Run()
	}

	go hub.Run()
}
