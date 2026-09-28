package complete

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"atomicptr.dev/peon/pkg/git"
	"github.com/urfave/cli/v3"
)

func WorktreeNameShellComplete(ctx context.Context, c *cli.Command) {
	if c.NArg() > 0 {
		return
	}

	g, err := git.FromContext(ctx)
	if err != nil {
		return
	}

	if !g.IsGitDir() {
		return
	}

	masterBranch := g.FindMasterBranch()
	if masterBranch != "" {
		parts := strings.Split(masterBranch, "/")
		fmt.Println(parts[len(parts)-1])
	}

	worktrees, err := g.GetWorktrees()
	if err != nil {
		return
	}

	commonRoot, err := g.FindCommonRoot()
	if err != nil {
		return
	}

	commonRootName := filepath.Base(commonRoot)

	for _, wt := range worktrees {
		pathName := filepath.Base(wt.Path)

		// skip path suffix matching on root dir
		if pathName == commonRootName {
			continue
		}

		// not managed by us, skip
		if !strings.HasPrefix(pathName, commonRootName) {
			continue
		}

		fmt.Println(wt.BranchShortName())
	}
}
