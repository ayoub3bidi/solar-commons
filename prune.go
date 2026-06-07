package solarcommons

import (
	"fmt"
	"os"
	"strings"
)

// PruneFeature removes paths and filters lines declared in a disabled Feature.
func PruneFeature(dest string, f Feature) error {
	for _, p := range f.RemovePaths {
		path, err := safeJoin(dest, p)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove feature path %q: %w", p, err)
		}
	}
	for _, lr := range f.RemoveLines {
		path, err := safeJoin(dest, lr.File)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read line removal file %q: %w", lr.File, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat line removal file %q: %w", lr.File, err)
		}
		lines := strings.Split(string(data), "\n")
		out := make([]string, 0, len(lines))
		for _, line := range lines {
			if !strings.Contains(line, lr.Pattern) {
				out = append(out, line)
			}
		}
		if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), info.Mode().Perm()); err != nil {
			return fmt.Errorf("write line removal file %q: %w", lr.File, err)
		}
	}
	return nil
}
