package git

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"atomicptr.dev/bits"
	"atomicptr.dev/deeperr"
	"atomicptr.dev/peon/pkg/constants"
)

type FlagKind int

type Worktree struct {
	Path       string
	Head       string
	Branch     string
	IsBare     bool
	IsDetached bool
	Locked     string
	Prunable   string
}

func (wt Worktree) BranchShortName() string {
	branch := wt.Branch

	if branch == "" {
		return ""
	}

	branch = strings.TrimPrefix(branch, "refs/heads/")
	branch = strings.TrimPrefix(branch, "refs/remotes/")
	branch = strings.TrimPrefix(branch, "refs/tags/")

	return branch
}

func (h Handler) GetWorktrees() ([]Worktree, error) {
	res, err := h.exec("worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}

	return h.parseWorktreeResults(res, "\x00")
}

func (h Handler) parseWorktreeResults(data, sep string) ([]Worktree, error) {
	var worktrees []Worktree
	curr := 0

	for entry := range strings.SplitSeq(data, sep) {
		if len(entry) == 0 {
			curr++
			continue
		}

		parts := strings.SplitN(entry, " ", 2)

		key := parts[0]

		switch key {
		case "worktree":
			path, err := filepath.Abs(parts[1])
			if err != nil {
				return nil, fmt.Errorf("could not parse worktree path: %w", err)
			}

			worktrees = append(worktrees, Worktree{Path: path})
		case "HEAD":
			worktrees[curr].Head = parts[1]
		case "branch":
			worktrees[curr].Branch = parts[1]
		case "bare":
			worktrees[curr].IsBare = true
		case "detached":
			worktrees[curr].IsDetached = true
		case "locked":
			worktrees[curr].Locked = parts[1]
		case "prunable":
			worktrees[curr].Prunable = parts[1]
		default:
			slog.Debug("unknown worktree key", "key", key, "path", h.WorkingDir)
		}
	}

	return worktrees, nil
}

func (h Handler) FindWorktreeByName(name string) (*Worktree, error) {
	worktrees, err := h.GetWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		branch := wt.BranchShortName()

		// has no branch name cant really be compared
		if branch == "" {
			continue
		}

		if strings.EqualFold(branch, name) {
			return &wt, nil
		}
	}

	// we couldnt find by branch, lets try finding by suffix
	commonRootDir, err := h.FindCommonRoot()
	if err != nil {
		return nil, err
	}

	commonRootName, err := filepath.Abs(commonRootDir)
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		wtDirName, err := filepath.Abs(wt.Path)
		if err != nil {
			return nil, err
		}

		// not managed by us (presumably)
		if !strings.HasPrefix(wtDirName, commonRootName) {
			continue
		}

		suffix := strings.TrimLeft(wtDirName[len(commonRootName):], "-")

		// if it matches the suffix, return it
		if strings.EqualFold(suffix, name) {
			return &wt, nil
		}
	}

	return nil, deeperr.NewWithCode(constants.ErrWorktreeNotFound, fmt.Sprintf("Could not find any worktree named: `%s`", name), nil)
}

func (h Handler) FindWorktreeByPath(path string) (*Worktree, error) {
	if !bits.PathExists(path) {
		return nil, fmt.Errorf("path `%s` does not exist", path)
	}

	searchPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	worktrees, err := h.GetWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		wtPath, err := filepath.Abs(wt.Path)
		if err != nil {
			return nil, err
		}

		if searchPath == wtPath {
			return &wt, nil
		}
	}

	return nil, fmt.Errorf("could not find worktree at the path `%s`", searchPath)
}

func (h Handler) CreateNewWorktreeFromMaster(newBranchName, targetPath string) (*Worktree, error) {
	branch := h.FindMasterBranch()
	if branch == "" {
		return nil, fmt.Errorf("could not determine branch name to fork off from")
	}

	return h.CreateNewWorktree(newBranchName, branch, targetPath)
}

func (h Handler) CreateNewWorktree(newBranchName, sourceBranch, targetPath string) (*Worktree, error) {
	if bits.PathExists(targetPath) {
		return nil, fmt.Errorf("create new worktree - target path `%s` already exists", targetPath)
	}

	if !h.BranchExists(sourceBranch) {
		return nil, fmt.Errorf("create new worktree - source branch `%s` does not exist", sourceBranch)
	}

	if h.BranchExists(newBranchName) {
		return nil, fmt.Errorf("create new worktree - branch `%s` already exists", newBranchName)
	}

	_, err := h.exec("worktree", "add", "-b", newBranchName, targetPath, sourceBranch)
	if err != nil {
		return nil, err
	}

	worktrees, err := h.GetWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		wtPath, err := filepath.Abs(wt.Path)
		if err != nil {
			slog.Error("create new worktree - could not get absolute path of worktree", "err", err)
			continue
		}

		if wtPath == targetPath {
			return &wt, nil
		}
	}

	return nil, fmt.Errorf("could not find newly created worktree")
}

func (h Handler) CreateWorktreeFromExistingBranch(branch, targetPath string) (*Worktree, error) {
	if !h.BranchExists(branch) {
		return nil, fmt.Errorf("branch `%s` does not exist", branch)
	}

	if bits.PathExists(targetPath) {
		return nil, fmt.Errorf("target path `%s` already exists", targetPath)
	}

	_, err := h.exec("worktree", "add", targetPath, branch)
	if err != nil {
		return nil, err
	}

	worktrees, err := h.GetWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		wtPath, err := filepath.Abs(wt.Path)
		if err != nil {
			slog.Error("could not get absolute path of worktree", "err", err)
			continue
		}

		if wtPath == targetPath {
			return &wt, nil
		}
	}

	return nil, fmt.Errorf("could not find newly created worktree")
}

func (h Handler) DeleteWorktree(wt *Worktree, force bool) error {
	// already deleted? Quit
	if !bits.PathExists(wt.Path) {
		return nil
	}

	// TODO: check for unstaged/unpushed changes

	var err error

	if force {
		_, err = h.exec("worktree", "remove", "-f", wt.Path)
	} else {
		_, err = h.exec("worktree", "remove", wt.Path)
	}

	return err
}
