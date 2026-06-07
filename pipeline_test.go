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
	dest := t.TempDir()

	// 1. Clone Mars at the pinned tag
	if err := CloneTemplate("https://github.com/ayoub3bidi/mars", "v0.1.1", dest); err != nil {
		t.Fatalf("CloneTemplate: %v", err)
	}
	if err := StripGit(dest); err != nil {
		t.Fatalf("StripGit: %v", err)
	}

	// 2. Parse manifest
	m, err := ParseManifest(dest)
	if err != nil {
		t.Fatalf("ParseManifest: %v", err)
	}
	if m.ID != "mars" {
		t.Fatalf("unexpected manifest ID: %q", m.ID)
	}

	// 3. Apply rewrites
	vars := map[string]string{
		"project_name": "my-api",
		"app_title":    "My API",
		"api_version":  "v1",
	}
	if err := ApplyRewrites(dest, m, vars); err != nil {
		t.Fatalf("ApplyRewrites: %v", err)
	}

	// 4. Prune OIDC feature
	var oidcFeature Feature
	for _, f := range m.Features {
		if f.Name == "oidc" {
			oidcFeature = f
			break
		}
	}
	if err := PruneFeature(dest, oidcFeature); err != nil {
		t.Fatalf("PruneFeature: %v", err)
	}

	// 5. Git init
	if err := InitGit(dest, "Initial commit"); err != nil {
		t.Fatalf("InitGit: %v", err)
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

	// OIDC path must be gone
	if _, err := os.Stat(filepath.Join(dest, "src", "routes", "oidc")); !os.IsNotExist(err) {
		t.Error("src/routes/oidc should have been removed")
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
	// Ensure there's no parent (single commit)
	if len(commit.ParentHashes) != 0 {
		t.Errorf("expected 0 parents, got %d", len(commit.ParentHashes))
	}
}
