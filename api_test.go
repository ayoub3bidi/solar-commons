package solarcommons

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/testutil"
)

func TestPublicFacade(t *testing.T) {
	source := testutil.NewScaffoldFixtureRepo(t)
	dest := filepath.Join(t.TempDir(), "my-api")

	result, err := Scaffold(ScaffoldConfig{
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
}

func TestPublicFacadeForceFlag(t *testing.T) {
	source := testutil.NewScaffoldFixtureRepo(t)
	dest := filepath.Join(t.TempDir(), "my-api")

	// First scaffold into dest
	_, err := Scaffold(ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Second scaffold into same dest without Force must fail
	_, err = Scaffold(ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
	})
	if !errors.Is(err, ErrDestinationExists) {
		t.Fatalf("expected ErrDestinationExists, got %v", err)
	}

	// Third scaffold with Force must succeed
	_, err = Scaffold(ScaffoldConfig{
		Repo:        source,
		Ref:         "v0.1.0",
		Destination: dest,
		Variables:   map[string]string{"project_name": "my-api", "app_title": "My API"},
		Force:       true,
	})
	if err != nil {
		t.Fatalf("Force scaffold: %v", err)
	}
}

func TestPublicFacadePreservesErrorCategories(t *testing.T) {
	_, err := Scaffold(ScaffoldConfig{})
	if err == nil {
		t.Fatal("expected invalid config error")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("error: got %v, want ErrInvalidConfig", err)
	}
}
