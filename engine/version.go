package engine

import (
	"fmt"
	"os"
	"runtime"
)

var (
	BuildID = "unknown" // git version id, passed in from shell
)

const (
	VERSION = "v0.0.1.beta"
	AUTHOR  = "funky.gao@gmail.com"
)

func ShowVersionAndExit() {
	fmt.Fprintf(os.Stderr, "%s %s (build: %s)\n", os.Args[0], VERSION,
		BuildID)
	fmt.Fprintf(os.Stderr, "Built with %s %s for %s/%s\n",
		runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	os.Exit(0)
}
