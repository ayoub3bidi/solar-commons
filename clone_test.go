package solarcommons

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestStripGit(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := StripGit(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); !os.IsNotExist(err) {
		t.Fatal(".git should not exist after StripGit")
	}
}

func TestCloneTemplate(t *testing.T) {
	source := newTaggedFixtureRepo(t)
	dest := t.TempDir()

	if err := CloneTemplate(source, "v0.1.0", dest); err != nil {
		t.Fatalf("CloneTemplate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); os.IsNotExist(err) {
		t.Fatal(".git should exist after clone")
	}
	// After strip, .git is gone.
	if err := StripGit(dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); !os.IsNotExist(err) {
		t.Fatal(".git should be gone after StripGit")
	}
}

func newTaggedFixtureRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatal(err)
	}
	hash, err := worktree.Commit("Initial commit", defaultTestCommitOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateTag("v0.1.0", hash, nil); err != nil {
		t.Fatal(err)
	}
	return dir
}

func defaultTestCommitOptions() *git.CommitOptions {
	return &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com"},
	}
}
