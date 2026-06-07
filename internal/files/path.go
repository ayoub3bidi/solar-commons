package files

import (
	"fmt"
	"path/filepath"
	"strings"
)

func SafeJoin(root, child string) (string, error) {
	if child == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	if filepath.IsAbs(child) {
		return "", fmt.Errorf("path %q must be relative", child)
	}

	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve root path: %w", err)
	}
	joined := filepath.Join(cleanRoot, child)
	cleanJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", child, err)
	}

	rel, err := filepath.Rel(cleanRoot, cleanJoined)
	if err != nil {
		return "", fmt.Errorf("compare path %q with root: %w", child, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes destination root", child)
	}

	return cleanJoined, nil
}
