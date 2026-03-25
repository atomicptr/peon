package fs

import "path/filepath"

func ValidatePath(path string) (string, bool) {
	if path == "" {
		return "", false
	}

	p, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}

	if Exists(p) {
		return p, true
	}

	return "", false
}
