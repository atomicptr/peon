package git

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atomicptr/peon/pkg/fs"
)

type contextKey string

const ContextKey contextKey = "github.com/atomicptr/peon/pkg/git"

type Handler struct {
	Executable string
	WorkingDir string
}

func (h Handler) exec(params ...string) (string, error) {
	cwd := h.WorkingDir

	if cwd == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not find working dir: %w", err)
		}

		cwd = wd
	}

	if !fs.Exists(cwd) {
		return "", fmt.Errorf("dir %s does not exist", cwd)
	}

	// #nosec G204 - this is fine
	cmd := exec.Command(filepath.Clean(h.Executable), params...)
	cmd.Dir = cwd

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Debug(
			"command failed",
			"git", h.Executable,
			"params", params,
			"err", err,
			"stdout", stdout.String(),
			"stderr", stderr.String(),
		)
		return "", fmt.Errorf("command `%s %s` failed: %w", h.Executable, params, err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func FromEnv() (*Handler, error) {
	path, ok := fs.ValidatePath(os.Getenv("PEON_GIT"))
	if ok {
		return &Handler{Executable: path}, nil
	}

	path, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("could not find git executable: %w", err)
	}

	return &Handler{Executable: path}, nil
}

func FromContext(ctx context.Context) (*Handler, error) {
	handler, ok := ctx.Value(ContextKey).(*Handler)
	if !ok {
		return nil, fmt.Errorf("git is not initialized")
	}

	return handler, nil
}
