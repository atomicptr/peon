package switchcmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/atomicptr/peon/pkg/cli/commands/shell"
	"github.com/atomicptr/peon/pkg/cli/common/orders"
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
	g, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !g.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", g.WorkingDir)
	}

	orders, hasOrders := orders.FromEnv()
	if !hasOrders {
		slog.Error(fmt.Sprintf("Shell integration not installed, please add `%s` to your `%s` file.", shell.EvalCommand(), shell.ConfigFile()))
	}

	name := cmd.StringArg("name")

	// no name specified, switch to default
	if name == "" {
		if hasOrders {
			rootDir, err := g.FindCommonRoot()
			if err != nil {
				return err
			}

			err = orders.ChangeDir(rootDir)
			if err != nil {
				return fmt.Errorf("could not change dirs: %w", err)
			}
		}

		return nil
	}

	// TODO: check if there is a worktree with the same literal name
	// TODO: check if there is a worktree with the same sluggified name

	if cmd.Bool("create") {
		// TODO: check if name conflicts, return error
		// TODO: create the worktree
		return nil // TODO: fall through cuz we switchin
	}

	// TODO: if name doesnt exist (and we havent created) fuzzy find the closest name (take recency into account)
	// TODO: if does not exist show error
	// TODO: switch

	return nil
}
