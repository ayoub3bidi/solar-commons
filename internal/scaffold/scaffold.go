package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
	gitops "github.com/ayoub3bidi/solar-commons/internal/git"
	"github.com/ayoub3bidi/solar-commons/internal/manifest"
	"github.com/ayoub3bidi/solar-commons/internal/prune"
	"github.com/ayoub3bidi/solar-commons/internal/rewrite"
	"github.com/ayoub3bidi/solar-commons/internal/types"
)

const (
	defaultInitialCommit = "Initial commit"
	defaultAuthorName    = "sun"
	defaultAuthorEmail   = "sun@solar"
)

// Run executes the full template generation pipeline.
func Run(config types.ScaffoldConfig) (types.ScaffoldResult, error) {
	config = normalizeConfig(config)
	if err := validateConfig(config); err != nil {
		return types.ScaffoldResult{}, err
	}
	if _, err := os.Stat(config.Destination); err == nil {
		if !config.Force {
			return types.ScaffoldResult{}, fmt.Errorf("%w: %s", apperrors.ErrDestinationExists, config.Destination)
		}
		if err := os.RemoveAll(config.Destination); err != nil {
			return types.ScaffoldResult{}, fmt.Errorf("remove existing destination %s: %w", config.Destination, err)
		}
	} else if !os.IsNotExist(err) {
		return types.ScaffoldResult{}, fmt.Errorf("check destination %s: %w", config.Destination, err)
	}

	if err := gitops.CloneTemplate(config.Repo, config.Ref, config.Destination); err != nil {
		return types.ScaffoldResult{}, err
	}
	if err := gitops.StripGit(config.Destination); err != nil {
		return types.ScaffoldResult{}, err
	}

	m, err := manifest.Parse(config.Destination)
	if err != nil {
		return types.ScaffoldResult{}, err
	}
	manifestPath := filepath.Join(config.Destination, "solar.manifest.yaml")
	if err := os.Remove(manifestPath); err != nil && !os.IsNotExist(err) {
		return types.ScaffoldResult{}, fmt.Errorf("remove template manifest: %w", err)
	}
	if err := rewrite.Apply(config.Destination, m.Rewrites, config.Variables); err != nil {
		return types.ScaffoldResult{}, err
	}

	prunedFeatures := make([]string, 0, len(config.DisabledFeatures))
	for _, featureName := range config.DisabledFeatures {
		feature, ok := m.FeatureByName(featureName)
		if !ok {
			return types.ScaffoldResult{}, fmt.Errorf("%w: %q", apperrors.ErrUnknownFeature, featureName)
		}
		if err := prune.Apply(config.Destination, feature); err != nil {
			return types.ScaffoldResult{}, err
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
		return types.ScaffoldResult{}, err
	}

	return types.ScaffoldResult{
		Destination:    config.Destination,
		TemplateID:     m.ID,
		Ref:            config.Ref,
		PrunedFeatures: prunedFeatures,
		InitialCommit:  config.InitialCommit,
	}, nil
}

func normalizeConfig(config types.ScaffoldConfig) types.ScaffoldConfig {
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

func validateConfig(config types.ScaffoldConfig) error {
	if config.Repo == "" {
		return fmt.Errorf("%w: repo must not be empty", apperrors.ErrInvalidConfig)
	}
	if config.Ref == "" {
		return fmt.Errorf("%w: ref must not be empty", apperrors.ErrInvalidConfig)
	}
	if config.Destination == "" {
		return fmt.Errorf("%w: destination must not be empty", apperrors.ErrInvalidConfig)
	}
	return nil
}
