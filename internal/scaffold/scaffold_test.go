package scaffold

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/apperrors"
	"github.com/ayoub3bidi/solar-commons/internal/testutil"
	"github.com/ayoub3bidi/solar-commons/internal/types"
	"github.com/go-git/go-git/v5"
)

func TestRun(t *testing.T) {
	source := testutil.NewScaffoldFixtureRepo(t)
	dest := filepath.Join(t.TempDir(), "my-api")

	result, err := Run(types.ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-api",
			"app_title":    "My API",
		},
		DisabledFeatures: []string{"oidc"},
		InitialCommit:    "Initial scaffold",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TemplateID != "fixture" {
		t.Fatalf("TemplateID: got %q, want %q", result.TemplateID, "fixture")
	}
	if len(result.PrunedFeatures) != 1 || result.PrunedFeatures[0] != "oidc" {
		t.Fatalf("PrunedFeatures: unexpected %+v", result.PrunedFeatures)
	}
	if _, err := os.Stat(filepath.Join(dest, "src", "routes", "oidc")); !os.IsNotExist(err) {
		t.Fatal("OIDC route directory should be pruned")
	}

	env, err := os.ReadFile(filepath.Join(dest, ".env.dist"))
	if err != nil {
		t.Fatal(err)
	}
	envContent := string(env)
	if strings.Contains(envContent, "fixture") || strings.Contains(envContent, "Fixture API") {
		t.Fatalf(".env.dist still contains fixture values:\n%s", envContent)
	}
	if strings.Contains(envContent, "OIDC") {
		t.Fatalf(".env.dist still contains OIDC values:\n%s", envContent)
	}
	if !strings.Contains(envContent, "my-api") || !strings.Contains(envContent, "My API") {
		t.Fatalf(".env.dist missing rewritten values:\n%s", envContent)
	}

	compose, err := os.ReadFile(filepath.Join(dest, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	composeContent := string(compose)
	if strings.Contains(composeContent, "fixture_api") || strings.Contains(composeContent, "fixture_db") {
		t.Fatalf("docker-compose.yml still contains fixture service names:\n%s", composeContent)
	}
	if !strings.Contains(composeContent, "my-api_api") || !strings.Contains(composeContent, "my-api_db") {
		t.Fatalf("docker-compose.yml missing rewritten service names:\n%s", composeContent)
	}

	resources, err := os.ReadFile(filepath.Join(dest, "src", "restful_ressources.py"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(resources), "oidc") {
		t.Fatalf("src/restful_ressources.py still contains OIDC references:\n%s", string(resources))
	}

	if _, err := os.Stat(filepath.Join(dest, ".git")); err != nil {
		t.Fatal(".git should exist after Run")
	}

	repo, err := git.PlainOpen(dest)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if commit.Message != "Initial scaffold" {
		t.Fatalf("commit message: got %q", commit.Message)
	}
	if len(commit.ParentHashes) != 0 {
		t.Fatalf("expected one clean initial commit, got %d parents", len(commit.ParentHashes))
	}
}

func TestRunForceFlag(t *testing.T) {
	source := testutil.NewScaffoldFixtureRepo(t)
	dest := filepath.Join(t.TempDir(), "my-api")

	// First scaffold
	_, err := Run(types.ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Without Force should fail
	_, err = Run(types.ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
	})
	if !errors.Is(err, apperrors.ErrDestinationExists) {
		t.Fatalf("expected ErrDestinationExists, got %v", err)
	}

	// With Force should succeed
	_, err = Run(types.ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
		Force:       true,
	})
	if err != nil {
		t.Fatalf("Force: %v", err)
	}
}

func TestRunRejectsExistingDestination(t *testing.T) {
	dest := t.TempDir()
	_, err := Run(types.ScaffoldConfig{
		Repo:        "unused",
		Ref:         "v0.1.0",
		Destination: dest,
	})
	if err == nil {
		t.Fatal("expected destination exists error")
	}
	if !errors.Is(err, apperrors.ErrDestinationExists) {
		t.Fatalf("error: got %v, want ErrDestinationExists", err)
	}
}

func TestRunErrorCategories(t *testing.T) {
	tests := map[string]struct {
		config types.ScaffoldConfig
		target error
	}{
		"invalid config": {
			config: types.ScaffoldConfig{},
			target: apperrors.ErrInvalidConfig,
		},
		"unknown feature": {
			config: types.ScaffoldConfig{
				Repo:             testutil.NewScaffoldFixtureRepo(t),
				Ref:              "v0.1.0",
				Destination:      filepath.Join(t.TempDir(), "unknown-feature"),
				Variables:        map[string]string{"project_name": "my-api", "app_title": "My API"},
				DisabledFeatures: []string{"missing"},
				InitialCommit:    "Initial scaffold",
			},
			target: apperrors.ErrUnknownFeature,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Run(tt.config)
			if err == nil {
				t.Fatal("expected error")
			}
			if !errors.Is(err, tt.target) {
				t.Fatalf("error: got %v, want %v", err, tt.target)
			}
		})
	}
}
