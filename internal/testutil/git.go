package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func NewTaggedRepo(t *testing.T) string {
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
	hash, err := worktree.Commit("Initial commit", DefaultCommitOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateTag("v0.1.0", hash, nil); err != nil {
		t.Fatal(err)
	}
	return dir
}

func NewScaffoldFixtureRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	CopyFixtureTemplate(t, dir)

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
	hash, err := worktree.Commit("Initial fixture", DefaultCommitOptions())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateTag("v0.1.0", hash, nil); err != nil {
		t.Fatal(err)
	}
	return dir
}

func DefaultCommitOptions() *git.CommitOptions {
	return &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com"},
	}
}
