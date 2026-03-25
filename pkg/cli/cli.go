package cli

import (
	"context"
	"os"

	"github.com/atomicptr/peon/pkg/cli/commands/list"
	"github.com/atomicptr/peon/pkg/cli/commands/remove"
	switchcmd "github.com/atomicptr/peon/pkg/cli/commands/switch"
	"github.com/urfave/cli/v3"
)

func Run() error {
	cmd := &cli.Command{
		Name:  "peon",
		Usage: "Peon is a Git Worktree management tool designed for humans",
		Commands: []*cli.Command{
			switchcmd.Command(),
			list.Command(),
			remove.Command(),
		},
	}

	return cmd.Run(context.Background(), os.Args)
}
