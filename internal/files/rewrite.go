package files

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
)

var templateVarPattern = regexp.MustCompile(`\{\{\.[A-Za-z_][A-Za-z0-9_]*\}\}`)

type Rewrite struct {
	File         string
	Replacements []Replacement
}

type Replacement struct {
	From string
	To   string
}

func ApplyRewrites(dest string, rewrites []Rewrite, vars map[string]string) error {
	for _, rewrite := range rewrites {
		path, err := SafeJoin(dest, rewrite.File)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read rewrite file %q: %w", rewrite.File, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat rewrite file %q: %w", rewrite.File, err)
		}
		content := string(data)
		for _, replacement := range rewrite.Replacements {
			to, err := ResolveVars(replacement.To, vars)
			if err != nil {
				return fmt.Errorf("resolve replacement for %q: %w", rewrite.File, err)
			}
			content = strings.ReplaceAll(content, replacement.From, to)
		}
		if err := os.WriteFile(path, []byte(content), info.Mode().Perm()); err != nil {
			return fmt.Errorf("write rewrite file %q: %w", rewrite.File, err)
		}
	}
	return nil
}

func ResolveVars(s string, vars map[string]string) (string, error) {
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
		return "", fmt.Errorf("%w: missing variable values: %s", apperrors.ErrMissingVariable, strings.Join(missing, ", "))
	}
	return out, nil
}
