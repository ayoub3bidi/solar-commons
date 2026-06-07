package files

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
)

func SafeJoin(root, child string) (string, error) {
	if child == "" {
		return "", fmt.Errorf("%w: path must not be empty", apperrors.ErrUnsafePath)
	}
	if filepath.IsAbs(child) {
		return "", fmt.Errorf("%w: path %q must be relative", apperrors.ErrUnsafePath, child)
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
		return "", fmt.Errorf("%w: path %q escapes destination root", apperrors.ErrUnsafePath, child)
	}

	return cleanJoined, nil
}
