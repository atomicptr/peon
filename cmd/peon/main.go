package main

import (
	"log/slog"
	"os"

	"atomicptr.dev/peon/pkg/cli"
)

func main() {
	err := cli.Run()
	if err != nil {
		slog.Error("caught error", "err", err)
		os.Exit(1)
	}
}
