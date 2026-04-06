package remove

import (
	"context"
	"fmt"

	"atomicptr.dev/deeperr"
	"atomicptr.dev/peon/pkg/cli/common/orders"
	"atomicptr.dev/peon/pkg/cli/common/usererr"
	"atomicptr.dev/peon/pkg/constants"
	"atomicptr.dev/peon/pkg/git"
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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Remove even if the worktree contains unstaged changes",
			},
		},
		Action: removeCommand,
	}
}

func removeCommand(ctx context.Context, cmd *cli.Command) error {
	g, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !g.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", g.WorkingDir)
	}

	orders, hasOrders := orders.FromEnv()
	if !hasOrders {
		usererr.ShellIntegrationNotInstalled()
	}

	rootDir, err := g.FindCommonRoot()
	if err != nil {
		return err
	}

	name := cmd.StringArg("name")
	forceRemove := cmd.Bool("force")

	// no name specified, switch to default
	if name == "" {
		if hasOrders {
			projectRoot, err := g.FindProjectRoot()
			if err != nil {
				return err
			}

			wt, err := g.FindWorktreeByPath(projectRoot)
			if err != nil {
				return err
			}

			if !forceRemove {
				untrackedChanges, err := g.StatusFor(projectRoot)
				if err != nil {
					return err
				}

				if len(untrackedChanges) > 0 {
					return deeperr.NewWithCode(constants.ErrUnstagedChanges, fmt.Sprintf("Repository `%s` contains unstaged changes, please commit them or run this command again with the `--force` flag", projectRoot), nil)
				}
			}

			err = g.DeleteWorktree(wt, forceRemove)
			if err != nil {
				return fmt.Errorf("couldnt remove worktree: %w", err)
			}

			return orders.ChangeDir(rootDir)
		}

		return nil
	}

	wt, err := g.FindWorktreeByName(name)
	if err != nil {
		return err
	}

	if !forceRemove {
		untrackedChanges, err := g.StatusFor(wt.Path)
		if err != nil {
			return err
		}

		if len(untrackedChanges) > 0 {
			return deeperr.NewWithCode(constants.ErrUnstagedChanges, fmt.Sprintf("Repository `%s` contains unstaged changes, please commit them or run this command again with the `--force` flag", wt.Path), nil)
		}
	}

	err = g.DeleteWorktree(wt, forceRemove)
	if err != nil {
		return fmt.Errorf("couldnt remove worktree: %w", err)
	}

	if !hasOrders {
		return nil
	}

	return orders.ChangeDir(rootDir)
}
