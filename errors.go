package solarcommons

import "github.com/ayoub3bidi/solar-commons/internal/apperrors"

var (
	ErrInvalidConfig     = apperrors.ErrInvalidConfig
	ErrDestinationExists = apperrors.ErrDestinationExists
	ErrInvalidManifest   = apperrors.ErrInvalidManifest
	ErrMissingVariable   = apperrors.ErrMissingVariable
	ErrUnknownFeature    = apperrors.ErrUnknownFeature
	ErrUnsafePath        = apperrors.ErrUnsafePath
)
