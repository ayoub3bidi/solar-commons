package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/files"
)

// PruneFeature removes paths and filters lines declared in a disabled Feature.
func PruneFeature(dest string, f Feature) error {
	return files.Prune(dest, toInternalPruneFeature(f))
}

func toInternalPruneFeature(feature Feature) files.PruneFeature {
	removeLines := make([]files.LineRemoval, 0, len(feature.RemoveLines))
	for _, lineRemoval := range feature.RemoveLines {
		removeLines = append(removeLines, files.LineRemoval{
			File:    lineRemoval.File,
			Pattern: lineRemoval.Pattern,
		})
	}
	return files.PruneFeature{
		RemovePaths: feature.RemovePaths,
		RemoveLines: removeLines,
	}
}
