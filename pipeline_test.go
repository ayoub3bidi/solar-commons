//go:build integration

package solarcommons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestFullPipeline(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "my-api")

	result, err := Scaffold(ScaffoldConfig{
		Repo:        "https://github.com/ayoub3bidi/mars",
		Ref:         "v0.1.3",
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

	// --- Assertions ---

	// .env.dist must not contain "mars" (case-sensitive project-name references)
	envData, err := os.ReadFile(filepath.Join(dest, ".env.dist"))
	if err != nil {
		t.Fatal(err)
	}
	envStr := string(envData)
	if strings.Contains(envStr, "mars") {
		t.Errorf(".env.dist still contains 'mars':\n%s", envStr)
	}
	if !strings.Contains(envStr, "my-api") {
		t.Errorf(".env.dist missing 'my-api':\n%s", envStr)
	}

	// OIDC env vars must be gone
	if strings.Contains(envStr, "OIDC") || strings.Contains(envStr, "GOOGLE_") {
		t.Errorf(".env.dist still contains OIDC/GOOGLE_ vars:\n%s", envStr)
	}

	// docker-compose.yml must not contain "mars" service names
	composeData, err := os.ReadFile(filepath.Join(dest, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	composeStr := string(composeData)
	if strings.Contains(composeStr, "marsdb") || strings.Contains(composeStr, "mars_cache") {
		t.Errorf("docker-compose.yml still contains unrewritten mars service names:\n%s", composeStr)
	}
	if !strings.Contains(composeStr, "my-api") {
		t.Errorf("docker-compose.yml missing 'my-api':\n%s", composeStr)
	}

	// OIDC paths must be gone
	for _, path := range []string{
		filepath.Join("src", "routes", "oidc"),
		filepath.Join("src", "controllers", "oidc"),
		filepath.Join("src", "utils", "oidc.py"),
	} {
		if _, err := os.Stat(filepath.Join(dest, path)); !os.IsNotExist(err) {
			t.Errorf("%s should have been removed by OIDC prune", path)
		}
	}

	// Git repo must have exactly one commit
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
	if commit.Message != "Initial commit" {
		t.Errorf("commit message: got %q", commit.Message)
	}
	if len(commit.ParentHashes) != 0 {
		t.Errorf("expected 0 parents, got %d", len(commit.ParentHashes))
	}
}
