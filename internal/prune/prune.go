package prune

import (
	"fmt"
	"os"
	"strings"

	"github.com/ayoub3bidi/solar-commons/internal/fs"
	"github.com/ayoub3bidi/solar-commons/internal/types"
)

func Apply(dest string, feature types.Feature) error {
	for _, p := range feature.RemovePaths {
		path, err := fs.SafeJoin(dest, p)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove feature path %q: %w", p, err)
		}
	}
	for _, lineRemoval := range feature.RemoveLines {
		path, err := fs.SafeJoin(dest, lineRemoval.File)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read line removal file %q: %w", lineRemoval.File, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat line removal file %q: %w", lineRemoval.File, err)
		}
		lines := strings.Split(string(data), "\n")
		out := make([]string, 0, len(lines))
		for _, line := range lines {
			if !strings.Contains(line, lineRemoval.Pattern) {
				out = append(out, line)
			}
		}
		if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), info.Mode().Perm()); err != nil {
			return fmt.Errorf("write line removal file %q: %w", lineRemoval.File, err)
		}
	}
	return nil
}
