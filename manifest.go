package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/manifest"
)

// ParseManifest reads solar.manifest.yaml from dest and returns the parsed Manifest.
func ParseManifest(dest string) (Manifest, error) {
	return manifest.Parse(dest)
}
