package switchcmd

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"github.com/atomicptr/peon/pkg/cli/commands/shell"
	"github.com/atomicptr/peon/pkg/cli/common/orders"
	"github.com/atomicptr/peon/pkg/git"
	"github.com/gosimple/slug"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/samber/lo"
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

	nameSlug := slug.Make(name)

	worktrees, err := g.GetWorktrees()
	if err != nil {
		return err
	}

	var matchedWt *git.Worktree = nil

	// look if there is a direct hit/match
	for _, wt := range worktrees {
		branch := wt.BranchShortName()

		// check for the exact same name
		if branch == name || strings.EqualFold(branch, name) || branch == nameSlug || strings.EqualFold(branch, nameSlug) {
			matchedWt = &wt
			break
		}

		pathName := strings.ToLower(filepath.Base(wt.Path))

		if strings.HasSuffix(pathName, "-"+strings.ToLower(name)) || strings.HasSuffix(pathName, "-"+nameSlug) {
			matchedWt = &wt
			break
		}
	}

	createNew := cmd.Bool("create")

	if createNew {
		// TODO: check if name conflicts, return error
		// TODO: create the worktree
		return nil // TODO: fall through cuz we switchin
	}

	// no wt found yet, try to fuzzy find
	if matchedWt == nil && !createNew {
		matchedWt = findClosestMatch(nameSlug, worktrees)
	}

	// still not found?
	if matchedWt == nil {
		return fmt.Errorf("could not find any worktree named like: %s", name)
	}

	// we found the worktree, switch to it
	return orders.ChangeDir(matchedWt.Path)
}

func findClosestMatch(query string, worktrees []git.Worktree) *git.Worktree {
	if len(worktrees) == 0 {
		return nil
	}

	targets := lo.Map(worktrees, func(wt git.Worktree, _ int) string {
		return strings.ToLower(wt.BranchShortName() + " " + filepath.Base(wt.Path))
	})

	matches := fuzzy.RankFind(query, targets)

	if len(matches) == 0 {
		return nil
	}

	sort.Sort(matches)

	bestMatch := matches[0].OriginalIndex
	return &worktrees[bestMatch]
}
