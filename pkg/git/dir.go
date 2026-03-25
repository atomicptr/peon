package git

import (
	"fmt"

	"github.com/atomicptr/peon/pkg/fs"
)

func (h Handler) IsGitDir() bool {
	res, err := h.exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false // error usually means we're not in a git directory
	}

	return res == "true"
}

func (h Handler) FindProjectRoot() (string, error) {
	res, err := h.exec("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}

	res, ok := fs.ValidatePath(res)
	if !ok {
		return "", fmt.Errorf("path %s does not exist", res)
	}

	return res, nil
}

func (h Handler) FindGitDir() (string, error) {
	res, err := h.exec("rev-parse", "--git-dir")
	if err != nil {
		return "", err
	}

	res, ok := fs.ValidatePath(res)
	if !ok {
		return "", fmt.Errorf("path %s does not exist", res)
	}

	return res, nil
}
