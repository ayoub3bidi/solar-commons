package solarcommons

import (
	"fmt"
	"os"
	"time"

	gitops "github.com/ayoub3bidi/solar-commons/internal/git"
)

const (
	defaultInitialCommit = "Initial commit"
	defaultAuthorName    = "sun"
	defaultAuthorEmail   = "sun@solar"
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

// Scaffold runs the full template generation pipeline.
func Scaffold(config ScaffoldConfig) (ScaffoldResult, error) {
	config = normalizeScaffoldConfig(config)
	if err := validateScaffoldConfig(config); err != nil {
		return ScaffoldResult{}, err
	}
	if _, err := os.Stat(config.Destination); err == nil {
		return ScaffoldResult{}, fmt.Errorf("destination already exists: %s", config.Destination)
	} else if !os.IsNotExist(err) {
		return ScaffoldResult{}, fmt.Errorf("check destination %s: %w", config.Destination, err)
	}

	if err := CloneTemplate(config.Repo, config.Ref, config.Destination); err != nil {
		return ScaffoldResult{}, err
	}
	if err := StripGit(config.Destination); err != nil {
		return ScaffoldResult{}, err
	}

	manifest, err := ParseManifest(config.Destination)
	if err != nil {
		return ScaffoldResult{}, err
	}
	if err := ApplyRewrites(config.Destination, manifest, config.Variables); err != nil {
		return ScaffoldResult{}, err
	}

	prunedFeatures := make([]string, 0, len(config.DisabledFeatures))
	for _, featureName := range config.DisabledFeatures {
		feature, ok := manifest.FeatureByName(featureName)
		if !ok {
			return ScaffoldResult{}, fmt.Errorf("unknown feature %q", featureName)
		}
		if err := PruneFeature(config.Destination, feature); err != nil {
			return ScaffoldResult{}, err
		}
		prunedFeatures = append(prunedFeatures, featureName)
	}

	if err := gitops.InitGit(
		config.Destination,
		config.InitialCommit,
		config.Author.Name,
		config.Author.Email,
		time.Now(),
	); err != nil {
		return ScaffoldResult{}, err
	}

	return ScaffoldResult{
		Destination:    config.Destination,
		TemplateID:     manifest.ID,
		Ref:            config.Ref,
		PrunedFeatures: prunedFeatures,
		InitialCommit:  config.InitialCommit,
	}, nil
}

func normalizeScaffoldConfig(config ScaffoldConfig) ScaffoldConfig {
	if config.InitialCommit == "" {
		config.InitialCommit = defaultInitialCommit
	}
	if config.Author.Name == "" {
		config.Author.Name = defaultAuthorName
	}
	if config.Author.Email == "" {
		config.Author.Email = defaultAuthorEmail
	}
	return config
}

func validateScaffoldConfig(config ScaffoldConfig) error {
	if config.Repo == "" {
		return fmt.Errorf("repo must not be empty")
	}
	if config.Ref == "" {
		return fmt.Errorf("ref must not be empty")
	}
	if config.Destination == "" {
		return fmt.Errorf("destination must not be empty")
	}
	return nil
}
