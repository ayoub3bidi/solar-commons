package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/rewrite"
)

// ApplyRewrites applies all file rewrites from the manifest to the project at dest.
// vars resolves {{.var_name}} tokens in replacement To values.
func ApplyRewrites(dest string, m Manifest, vars map[string]string) error {
	return rewrite.Apply(dest, m.Rewrites, vars)
}
