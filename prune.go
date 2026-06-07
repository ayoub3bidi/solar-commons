package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/prune"
)

// PruneFeature removes paths and filters lines declared in a disabled Feature.
func PruneFeature(dest string, f Feature) error {
	return prune.Apply(dest, f)
}
