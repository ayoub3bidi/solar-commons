package solarcommons

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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

// ParseManifest reads solar.manifest.yaml from dest and returns the parsed Manifest.
func ParseManifest(dest string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(dest, "solar.manifest.yaml"))
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse manifest: %w", err)
	}
	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// Validate checks the required manifest contract before any scaffold operation runs.
func (m Manifest) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("manifest id must not be empty")
	}
	if m.Repo == "" {
		return fmt.Errorf("manifest repo must not be empty")
	}
	if m.DefaultRef == "" {
		return fmt.Errorf("manifest default_ref must not be empty")
	}

	seenVariables := make(map[string]struct{}, len(m.Variables))
	for _, variable := range m.Variables {
		if variable.Name == "" {
			return fmt.Errorf("manifest variable name must not be empty")
		}
		if _, ok := seenVariables[variable.Name]; ok {
			return fmt.Errorf("manifest variable %q is duplicated", variable.Name)
		}
		seenVariables[variable.Name] = struct{}{}
	}

	for _, rewrite := range m.Rewrites {
		if rewrite.File == "" {
			return fmt.Errorf("manifest rewrite file must not be empty")
		}
		for _, replacement := range rewrite.Replacements {
			if replacement.From == "" {
				return fmt.Errorf("manifest rewrite %q has an empty replacement source", rewrite.File)
			}
		}
	}

	seenFeatures := make(map[string]struct{}, len(m.Features))
	for _, feature := range m.Features {
		if feature.Name == "" {
			return fmt.Errorf("manifest feature name must not be empty")
		}
		if _, ok := seenFeatures[feature.Name]; ok {
			return fmt.Errorf("manifest feature %q is duplicated", feature.Name)
		}
		seenFeatures[feature.Name] = struct{}{}
		for _, removal := range feature.RemoveLines {
			if removal.File == "" {
				return fmt.Errorf("manifest feature %q has an empty remove_lines file", feature.Name)
			}
			if removal.Pattern == "" {
				return fmt.Errorf("manifest feature %q has an empty remove_lines pattern", feature.Name)
			}
		}
	}

	return nil
}
