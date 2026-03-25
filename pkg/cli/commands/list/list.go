package list

import (
	"context"
	"fmt"

	"github.com/atomicptr/peon/pkg/git"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all available worktrees",
		Aliases: []string{"ls"},
		Action:  listCommand,
	}
}

func listCommand(ctx context.Context, cmd *cli.Command) error {
	git, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !git.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", git.WorkingDir)
	}

	return nil
}
