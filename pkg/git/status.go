package git

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

type StatusKind int

const (
	statusInvalid StatusKind = iota
	StatusUntracked
	StatusModified
	StatusAdded
	StatusDeleted
	StatusRenamed
)

type FileStatus struct {
	Status StatusKind
	Path   string
}

func (h Handler) Status() ([]FileStatus, error) {
	return h.StatusFor(h.WorkingDir)
}

func (h Handler) StatusFor(path string) ([]FileStatus, error) {
	s, err := h.execIn(path, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}

	if len(s) == 0 {
		return nil, nil
	}

	lines := strings.Split(s, "\n")

	status := make([]FileStatus, 0, len(lines))

	for _, line := range lines {
		l := strings.TrimSpace(line)
		parts := strings.SplitN(l, " ", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid git status line found: `%s`", line)
		}

		statusKind, err := parseStatusKind(parts[0])
		if err != nil {
			return nil, err
		}

		statusPath, err := filepath.Abs(filepath.Join(path, parts[1]))
		if err != nil {
			return nil, err
		}

		status = append(status, FileStatus{Status: statusKind, Path: statusPath})
	}

	slog.Debug("???? YOLO", "status", status)

	return status, nil
}

func parseStatusKind(k string) (StatusKind, error) {
	switch k {
	case "??":
		return StatusUntracked, nil
	case "M":
		return StatusModified, nil
	case "A":
		return StatusAdded, nil
	case "D":
		return StatusDeleted, nil
	case "R":
		return StatusRenamed, nil
	}

	return statusInvalid, fmt.Errorf("unknown git status kind: %s", k)
}
