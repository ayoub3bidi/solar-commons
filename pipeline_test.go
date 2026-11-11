//go:build integration

package solarcommons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
)

func assertFileNoPlanet(t *testing.T, dest, rel, planet string) {
	t.Helper()
	path := filepath.Join(dest, rel)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if strings.Contains(string(data), planet) {
		t.Errorf("%s still contains %q", rel, planet)
	}
}

func assertSingleCommit(t *testing.T, dest, message string) {
	t.Helper()
	repo, err := git.PlainOpen(dest)
	if err != nil {
		t.Fatalf("PlainOpen: %v", err)
	}
	ref, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if commit.Message != message {
		t.Errorf("commit message: got %q", commit.Message)
	}
	if len(commit.ParentHashes) != 0 {
		t.Errorf("expected 0 parents, got %d", len(commit.ParentHashes))
	}
}

func assertAPIRuntimeRewritten(t *testing.T, dest, planet string) {
	t.Helper()
	for _, rel := range []string{
		".env.dist",
		"docker-compose.yml",
		"src/constants/settings.py",
	} {
		assertFileNoPlanet(t, dest, rel, planet)
	}
}

func TestFullPipelineMars(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "my-api")

	result, err := Scaffold(ScaffoldConfig{
		Repo:        "https://github.com/ayoub3bidi/mars",
		Ref:         "v0.1.5",
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-api",
			"app_title":    "My API",
			"api_version":  "v1",
		},
		DisabledFeatures: []string{"oidc"},
		InitialCommit:    "Initial commit",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "mars" {
		t.Fatalf("TemplateID: got %q, want %q", result.TemplateID, "mars")
	}

	envStr, err := os.ReadFile(filepath.Join(dest, ".env.dist"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(envStr), "OIDC") || strings.Contains(string(envStr), "GOOGLE_") {
		t.Errorf(".env.dist still contains OIDC vars")
	}

	for _, path := range []string{
		filepath.Join("src", "routes", "oidc"),
		filepath.Join("src", "controllers", "oidc"),
		filepath.Join("src", "utils", "oidc.py"),
	} {
		if _, err := os.Stat(filepath.Join(dest, path)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed", path)
		}
	}

	if _, err := os.Stat(filepath.Join(dest, "solar.manifest.yaml")); !os.IsNotExist(err) {
		t.Error("solar.manifest.yaml should be removed from scaffold output")
	}

	assertAPIRuntimeRewritten(t, dest, "mars")
	assertSingleCommit(t, dest, "Initial commit")
}

func TestFullPipelineMercury(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "my-flyway-api")

	result, err := Scaffold(ScaffoldConfig{
		Repo:        "https://github.com/ayoub3bidi/mercury",
		Ref:         "v0.4.2",
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-flyway-api",
			"app_title":    "My Flyway API",
			"api_version":  "v1",
		},
		DisabledFeatures: []string{"oidc"},
		InitialCommit:    "Initial commit",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "mercury" {
		t.Fatalf("TemplateID: got %q, want %q", result.TemplateID, "mercury")
	}

	assertAPIRuntimeRewritten(t, dest, "mercury")
	assertSingleCommit(t, dest, "Initial commit")
}

func TestFullPipelineVenus(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "my-docs")

	result, err := Scaffold(ScaffoldConfig{
		Repo:        "https://github.com/ayoub3bidi/venus",
		Ref:         "v1.1.2",
		Destination: dest,
		Variables: map[string]string{
			"project_name": "my-docs",
			"site_title":   "My Docs",
		},
		DisabledFeatures: []string{"i18n", "search"},
		InitialCommit:    "Initial commit",
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if result.TemplateID != "venus" {
		t.Fatalf("TemplateID: got %q, want %q", result.TemplateID, "venus")
	}

	if _, err := os.Stat(filepath.Join(dest, "i18n", "fr")); !os.IsNotExist(err) {
		t.Error("i18n/fr should be removed")
	}
	if _, err := os.Stat(filepath.Join(dest, "src", "plugins", "local-search.cjs")); !os.IsNotExist(err) {
		t.Error("local-search plugin should be removed when search disabled")
	}

	assertFileNoPlanet(t, dest, "package.json", "venus")
	assertFileNoPlanet(t, dest, "docusaurus.config.ts", "venus")
	assertSingleCommit(t, dest, "Initial commit")
}
