package remove

import (
	"context"
	"fmt"

	"github.com/atomicptr/peon/pkg/git"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove worktree",
		Aliases:   []string{"rm"},
		ArgsUsage: "[name]",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "name",
				UsageText: "Worktree name",
			},
		},
		Action: removeCommand,
	}
}

func removeCommand(ctx context.Context, cmd *cli.Command) error {
	git, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !git.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", git.WorkingDir)
	}

	return nil
}
