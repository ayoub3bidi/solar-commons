package solarcommons

import (
	"os"
	"path/filepath"
	"strings"
)

// PruneFeature removes paths and filters lines declared in a disabled Feature.
func PruneFeature(dest string, f Feature) error {
	for _, p := range f.RemovePaths {
		if err := os.RemoveAll(filepath.Join(dest, p)); err != nil {
			return err
		}
	}
	for _, lr := range f.RemoveLines {
		path := filepath.Join(dest, lr.File)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		out := make([]string, 0, len(lines))
		for _, line := range lines {
			if !strings.Contains(line, lr.Pattern) {
				out = append(out, line)
			}
		}
		if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644); err != nil {
			return err
		}
	}
	return nil
}
