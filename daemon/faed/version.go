package main

import (
	"fmt"
	"github.com/funkygao/fae/engine"
	"os"
	"runtime"
)

const (
	VERSION = "v0.0.1.alpha"
	AUTHOR  = "funky.gao@gmail.com"
)

func showVersionAndExit() {
	fmt.Fprintf(os.Stderr, "%s %s (build: %s)\n", os.Args[0], VERSION,
		engine.BuildID)
	fmt.Fprintf(os.Stderr, "Built with %s %s for %s/%s\n",
		runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}
