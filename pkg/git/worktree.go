package git

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
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
