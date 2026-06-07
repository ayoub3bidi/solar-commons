package types

import (
	"fmt"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
)

type ScaffoldConfig struct {
	Repo             string
	Ref              string
	Destination      string
	Variables        map[string]string
	DisabledFeatures []string
	InitialCommit    string
	Author           GitAuthor
}

type GitAuthor struct {
	Name  string
	Email string
}

type ScaffoldResult struct {
	Destination    string
	TemplateID     string
	Ref            string
	PrunedFeatures []string
	InitialCommit  string
}

type Manifest struct {
	ID         string        `yaml:"id"`
	Repo       string        `yaml:"repo"`
	DefaultRef string        `yaml:"default_ref"`
	Variables  []Variable    `yaml:"variables"`
	Rewrites   []FileRewrite `yaml:"rewrites"`
	Features   []Feature     `yaml:"features"`
}

type Variable struct {
	Name string `yaml:"name"`
}

type FileRewrite struct {
	File         string        `yaml:"file"`
	Replacements []Replacement `yaml:"replacements"`
}

type Replacement struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Feature struct {
	Name        string        `yaml:"name"`
	Default     bool          `yaml:"default"`
	RemovePaths []string      `yaml:"remove_paths"`
	RemoveLines []LineRemoval `yaml:"remove_lines"`
}

type LineRemoval struct {
	File    string `yaml:"file"`
	Pattern string `yaml:"pattern"`
}

// Validate checks the required manifest contract before any scaffold operation runs.
func (m Manifest) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("%w: manifest id must not be empty", apperrors.ErrInvalidManifest)
	}
	if m.Repo == "" {
		return fmt.Errorf("%w: manifest repo must not be empty", apperrors.ErrInvalidManifest)
	}
	if m.DefaultRef == "" {
		return fmt.Errorf("%w: manifest default_ref must not be empty", apperrors.ErrInvalidManifest)
	}

	seenVariables := make(map[string]struct{}, len(m.Variables))
	for _, variable := range m.Variables {
		if variable.Name == "" {
			return fmt.Errorf("%w: manifest variable name must not be empty", apperrors.ErrInvalidManifest)
		}
		if _, ok := seenVariables[variable.Name]; ok {
			return fmt.Errorf("%w: manifest variable %q is duplicated", apperrors.ErrInvalidManifest, variable.Name)
		}
		seenVariables[variable.Name] = struct{}{}
	}

	for _, rewrite := range m.Rewrites {
		if rewrite.File == "" {
			return fmt.Errorf("%w: manifest rewrite file must not be empty", apperrors.ErrInvalidManifest)
		}
		for _, replacement := range rewrite.Replacements {
			if replacement.From == "" {
				return fmt.Errorf("%w: manifest rewrite %q has an empty replacement source", apperrors.ErrInvalidManifest, rewrite.File)
			}
		}
	}

	seenFeatures := make(map[string]struct{}, len(m.Features))
	for _, feature := range m.Features {
		if feature.Name == "" {
			return fmt.Errorf("%w: manifest feature name must not be empty", apperrors.ErrInvalidManifest)
		}
		if _, ok := seenFeatures[feature.Name]; ok {
			return fmt.Errorf("%w: manifest feature %q is duplicated", apperrors.ErrInvalidManifest, feature.Name)
		}
		seenFeatures[feature.Name] = struct{}{}
		for _, removal := range feature.RemoveLines {
			if removal.File == "" {
				return fmt.Errorf("%w: manifest feature %q has an empty remove_lines file", apperrors.ErrInvalidManifest, feature.Name)
			}
			if removal.Pattern == "" {
				return fmt.Errorf("%w: manifest feature %q has an empty remove_lines pattern", apperrors.ErrInvalidManifest, feature.Name)
			}
		}
	}

	return nil
}

// FeatureByName returns the feature with name if it exists in the manifest.
func (m Manifest) FeatureByName(name string) (Feature, bool) {
	for _, feature := range m.Features {
		if feature.Name == name {
			return feature, true
		}
	}
	return Feature{}, false
}
