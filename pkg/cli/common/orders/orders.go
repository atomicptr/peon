package orders

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/atomicptr/peon/pkg/fs"
)

type Order struct {
	Path string
}

func (o Order) ChangeDir(path string) error {
	if !fs.IsDir(path) {
		return fmt.Errorf("path %s is not a directory or does not exist", path)
	}

	return os.WriteFile(o.Path, fmt.Appendf(nil, `__peon_cd "%s"`, path), 0600)
}

func FromEnv() (Order, bool) {
	orders := os.Getenv("PEON_ORDERS")
	if orders == "" {
		return Order{}, false
	}

	path, ok := fs.ValidatePath(strings.ReplaceAll(filepath.Clean(orders), "\n", ""))
	if !ok {
		slog.Debug("No orders have been supplied")
		return Order{}, false
	}

	slog.Debug("Received orders!", "path", path)

	return Order{
		Path: path,
	}, true
}
