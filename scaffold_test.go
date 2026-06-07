package solarcommons

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
	if len(commit.ParentHashes) != 0 {
		t.Fatalf("expected one clean initial commit, got %d parents", len(commit.ParentHashes))
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
	copyFixtureTemplate(t, dir)

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

func copyFixtureTemplate(t *testing.T, dest string) {
	t.Helper()

	source := filepath.Join("internal", "fixtures", "template")
	err := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relativePath == "." {
			return nil
		}
		target := filepath.Join(dest, relativePath)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
}
