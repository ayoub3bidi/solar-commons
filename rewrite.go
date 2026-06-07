package solarcommons

import (
	"os"
	"path/filepath"
	"strings"
)

// ApplyRewrites applies all file rewrites from the manifest to the project at dest.
// vars resolves {{.var_name}} tokens in replacement To values.
func ApplyRewrites(dest string, m Manifest, vars map[string]string) error {
	for _, rw := range m.Rewrites {
		path := filepath.Join(dest, rw.File)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		for _, r := range rw.Replacements {
			to := resolveVars(r.To, vars)
			content = strings.ReplaceAll(content, r.From, to)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// resolveVars replaces {{.key}} tokens in s with values from vars.
func resolveVars(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{."+k+"}}", v)
	}
	return s
}
