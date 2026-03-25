package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/atomicptr/peon/pkg/fs"
)

type contextKey string

const ContextKey contextKey = "github.com/atomicptr/peon/pkg/git"

type Handler struct {
	Executable string
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

func FromContext(ctx context.Context) *Handler {
	handler, ok := ctx.Value(ContextKey).(*Handler)
	if !ok {
		panic("git handler uninitialized")
	}

	return handler
}
