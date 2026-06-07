package solarcommons

import (
	"github.com/ayoub3bidi/solar-commons/internal/scaffold"
)

// Scaffold runs the full template generation pipeline.
func Scaffold(config ScaffoldConfig) (ScaffoldResult, error) {
	return scaffold.Run(config)
}
