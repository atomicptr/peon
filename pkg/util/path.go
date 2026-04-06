package util

import (
	"path/filepath"
	"strings"
)

func CutPathAfter(path string, marker string) string {
	p := filepath.ToSlash(path)

	_, after, ok := strings.Cut(p, marker)
	if !ok {
		return ""
	}

	return after
}
