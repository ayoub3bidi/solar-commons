package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/files"
)

// ApplyRewrites applies all file rewrites from the manifest to the project at dest.
// vars resolves {{.var_name}} tokens in replacement To values.
func ApplyRewrites(dest string, m Manifest, vars map[string]string) error {
	return files.ApplyRewrites(dest, toInternalRewrites(m.Rewrites), vars)
}

func toInternalRewrites(rewrites []FileRewrite) []files.Rewrite {
	out := make([]files.Rewrite, 0, len(rewrites))
	for _, rewrite := range rewrites {
		replacements := make([]files.Replacement, 0, len(rewrite.Replacements))
		for _, replacement := range rewrite.Replacements {
			replacements = append(replacements, files.Replacement{
				From: replacement.From,
				To:   replacement.To,
			})
		}
		out = append(out, files.Rewrite{
			File:         rewrite.File,
			Replacements: replacements,
		})
	}
	return out
}
