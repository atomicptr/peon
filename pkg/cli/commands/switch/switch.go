package switchcmd

import (
	"context"
	"fmt"

	"github.com/atomicptr/peon/pkg/git"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "Switch to a worktree, create a new one if needed",
		Aliases:   []string{"sw"},
		ArgsUsage: "[name]",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "name",
				UsageText: "Worktree name",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Create a new worktree",
			},
		},
		Action: switchCommand,
	}
}

func switchCommand(ctx context.Context, cmd *cli.Command) error {
	git, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !git.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", git.WorkingDir)
	}

	return nil
}
