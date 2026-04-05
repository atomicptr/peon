package fs

import (
	"path/filepath"

	"atomicptr.dev/bits"
)

func ValidatePath(path string) (string, bool) {
	if path == "" {
		return "", false
	}

	p, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	if bits.PathExists(p) {
		return p, true
	}

	return "", false
}
