package solarcommons

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var templateVarPattern = regexp.MustCompile(`\{\{\.[A-Za-z_][A-Za-z0-9_]*\}\}`)

// ApplyRewrites applies all file rewrites from the manifest to the project at dest.
// vars resolves {{.var_name}} tokens in replacement To values.
func ApplyRewrites(dest string, m Manifest, vars map[string]string) error {
	for _, rw := range m.Rewrites {
		path, err := safeJoin(dest, rw.File)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read rewrite file %q: %w", rw.File, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat rewrite file %q: %w", rw.File, err)
		}
		content := string(data)
		for _, r := range rw.Replacements {
			to, err := resolveVars(r.To, vars)
			if err != nil {
				return fmt.Errorf("resolve replacement for %q: %w", rw.File, err)
			}
			content = strings.ReplaceAll(content, r.From, to)
		}
		if err := os.WriteFile(path, []byte(content), info.Mode().Perm()); err != nil {
			return fmt.Errorf("write rewrite file %q: %w", rw.File, err)
		}
	}
	return nil
}

// resolveVars replaces {{.key}} tokens in s with values from vars.
func resolveVars(s string, vars map[string]string) (string, error) {
	var missing []string
	out := templateVarPattern.ReplaceAllStringFunc(s, func(token string) string {
		key := strings.TrimSuffix(strings.TrimPrefix(token, "{{."), "}}")
		value, ok := vars[key]
		if !ok {
			missing = append(missing, key)
			return token
		}
		return value
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("missing variable values: %s", strings.Join(missing, ", "))
	}
	return out, nil
}
