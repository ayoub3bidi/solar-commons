package apperrors

import "errors"

var (
	ErrInvalidConfig     = errors.New("invalid config")
	ErrDestinationExists = errors.New("destination exists")
	ErrInvalidManifest   = errors.New("invalid manifest")
	ErrMissingVariable   = errors.New("missing variable")
	ErrUnknownFeature    = errors.New("unknown feature")
	ErrUnsafePath        = errors.New("unsafe path")
)
