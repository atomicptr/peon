package main

import (
	"log/slog"
	"os"

	"github.com/atomicptr/peon/pkg/cli"
)

func main() {
	err := cli.Run()
	if err != nil {
		slog.Error("caught error", "err", err)
		os.Exit(1)
	}
}
