package types

import (
	"errors"
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
)

func TestManifestValidateRejectsInvalidManifests(t *testing.T) {
	tests := map[string]Manifest{
		"missing id": {
			Repo:       "https://example.com/repo",
			DefaultRef: "v0.1.0",
		},
		"duplicate variable": {
			ID:         "fixture",
			Repo:       "https://example.com/repo",
			DefaultRef: "v0.1.0",
			Variables:  []Variable{{Name: "project_name"}, {Name: "project_name"}},
		},
		"empty replacement source": {
			ID:         "fixture",
			Repo:       "https://example.com/repo",
			DefaultRef: "v0.1.0",
			Rewrites: []FileRewrite{{
				File:         ".env.dist",
				Replacements: []Replacement{{To: "value"}},
			}},
		},
		"empty line removal pattern": {
			ID:         "fixture",
			Repo:       "https://example.com/repo",
			DefaultRef: "v0.1.0",
			Features: []Feature{{
				Name:        "oidc",
				RemoveLines: []LineRemoval{{File: ".env.dist"}},
			}},
		},
	}

	for name, manifest := range tests {
		t.Run(name, func(t *testing.T) {
			if err := manifest.Validate(); err == nil {
				t.Fatal("expected validation error")
			} else if !errors.Is(err, apperrors.ErrInvalidManifest) {
				t.Fatalf("error: got %v, want ErrInvalidManifest", err)
			}
		})
	}
}
