package switchcmd

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"strings"

	"atomicptr.dev/deeperr"
	"atomicptr.dev/peon/pkg/cli/common/orders"
	"atomicptr.dev/peon/pkg/cli/common/usererr"
	"atomicptr.dev/peon/pkg/constants"
	"atomicptr.dev/peon/pkg/git"
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
		usererr.ShellIntegrationNotInstalled()
	}

	name := cmd.StringArg("name")

	// no name specified, switch to default
	if name == "" {
		if hasOrders {
			rootDir, err := g.FindCommonRoot()
			if err != nil {
				return err
			}

			return orders.ChangeDir(rootDir)
		}

		return nil
	}

	nameSlug := slug.Make(name)

	worktrees, err := g.GetWorktrees()
	if err != nil {
		return err
	}

	commonRoot, err := g.FindCommonRoot()
	if err != nil {
		return err
	}

	commonRootName := filepath.Base(commonRoot)

	var matchedWt *git.Worktree = nil

	// look if there is a direct hit/match
	for _, wt := range worktrees {
		branch := wt.BranchShortName()

		// check for the exact same name
		if strings.EqualFold(branch, name) || strings.EqualFold(branch, nameSlug) {
			matchedWt = &wt
			break
		}

		pathName := filepath.Base(wt.Path)

		// skip path suffix matching on root dir
		if pathName == commonRootName {
			continue
		}

		// not managed by us, skip
		if !strings.HasPrefix(pathName, commonRootName) {
			continue
		}

		suffix := strings.TrimLeft(pathName[len(commonRootName):], "-")

		if strings.EqualFold(suffix, name) {
			matchedWt = &wt
			break
		}
	}

	createNew := cmd.Bool("create")

	if createNew {
		if matchedWt != nil {
			return fmt.Errorf("can't create worktree `%s`, because it collides with `%s` (%s)", name, matchedWt.Path, matchedWt.BranchShortName())
		}

		targetPath := filepath.Join(filepath.Dir(commonRoot), fmt.Sprintf("%s-%s", commonRootName, nameSlug))

		slog.Debug("creating new worktree...", "path", targetPath)

		if g.BranchExists(name) {
			wt, err := g.CreateWorktreeFromExistingBranch(name, targetPath)
			if err != nil {
				return fmt.Errorf("could not create worktree from branch `%s`: %w", name, err)
			}

			matchedWt = wt
		} else {
			// TODO: also implement branching off a different branch then master
			wt, err := g.CreateNewWorktreeFromMaster(name, targetPath)
			if err != nil {
				return fmt.Errorf("could not create new worktree: %w", err)
			}

			matchedWt = wt
		}
	}

	// no wt found yet, try to fuzzy find
	if matchedWt == nil && !createNew {
		matchedWt = findClosestMatch(nameSlug, worktrees)
	}

	// still not found?
	if matchedWt == nil {
		return deeperr.NewWithCode(constants.ErrWorktreeNotFound, fmt.Sprintf("Could not find any worktree named: `%s`", name), nil)
	}

	if !hasOrders {
		return nil
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
