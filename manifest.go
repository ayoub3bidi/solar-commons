package solarcommons

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	ID         string       `yaml:"id"`
	Repo       string       `yaml:"repo"`
	DefaultRef string       `yaml:"default_ref"`
	Variables  []Variable   `yaml:"variables"`
	Rewrites   []FileRewrite `yaml:"rewrites"`
	Features   []Feature    `yaml:"features"`
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
		return Manifest{}, err
	}
	var m Manifest
	return m, yaml.Unmarshal(data, &m)
}
