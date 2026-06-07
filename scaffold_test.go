package solarcommons

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestScaffold(t *testing.T) {
	source := newScaffoldFixtureRepo(t)
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
	if len(result.PrunedFeatures) != 1 || result.PrunedFeatures[0] != "oidc" {
		t.Fatalf("PrunedFeatures: unexpected %+v", result.PrunedFeatures)
	}
	if _, err := os.Stat(filepath.Join(dest, "src", "routes", "oidc")); !os.IsNotExist(err) {
		t.Fatal("OIDC route directory should be pruned")
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); err != nil {
		t.Fatal(".git should exist after Scaffold")
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
}

func TestScaffoldRejectsExistingDestination(t *testing.T) {
	dest := t.TempDir()
	_, err := Scaffold(ScaffoldConfig{
		Repo:        "unused",
		Ref:         "v0.1.0",
		Destination: dest,
	})
	if err == nil {
		t.Fatal("expected destination exists error")
	}
	if !errors.Is(err, ErrDestinationExists) {
		t.Fatalf("error: got %v, want ErrDestinationExists", err)
	}
}

func TestScaffoldErrorCategories(t *testing.T) {
	tests := map[string]struct {
		config ScaffoldConfig
		target error
	}{
		"invalid config": {
			config: ScaffoldConfig{},
			target: ErrInvalidConfig,
		},
		"unknown feature": {
			config: ScaffoldConfig{
				Repo:             newScaffoldFixtureRepo(t),
				Ref:              "v0.1.0",
				Destination:      filepath.Join(t.TempDir(), "unknown-feature"),
				Variables:        map[string]string{"project_name": "my-api", "app_title": "My API"},
				DisabledFeatures: []string{"missing"},
				InitialCommit:    "Initial scaffold",
			},
			target: ErrUnknownFeature,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Scaffold(tt.config)
			if err == nil {
				t.Fatal("expected error")
			}
			if !errors.Is(err, tt.target) {
				t.Fatalf("error: got %v, want %v", err, tt.target)
			}
		})
	}
}

func newScaffoldFixtureRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	writeFile(t, dir, "solar.manifest.yaml", `id: fixture
repo: fixture
default_ref: v0.1.0
variables:
  - name: project_name
  - name: app_title
rewrites:
  - file: .env.dist
    replacements:
      - from: "fixture"
        to: "{{.project_name}}"
      - from: "Fixture API"
        to: "{{.app_title}}"
features:
  - name: oidc
    default: false
    remove_paths:
      - src/routes/oidc
    remove_lines:
      - file: .env.dist
        pattern: OIDC
`)
	writeFile(t, dir, ".env.dist", "APP_NAME=fixture\nAPP_TITLE=Fixture API\nOIDC_CLIENT_ID=changeme\n")
	writeFile(t, dir, "src/routes/oidc/google.py", "# oidc\n")

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if err := worktree.AddGlob("."); err != nil {
		t.Fatal(err)
	}
	hash, err := worktree.Commit("Initial fixture", defaultTestCommitOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateTag("v0.1.0", hash, nil); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeFile(t *testing.T, root, path, content string) {
	t.Helper()

	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
