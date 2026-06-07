package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
)

func TestInitGit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	const msg = "Initial commit"
	if err := InitGit(dir, msg, "sun", "sun@solar", time.Now()); err != nil {
		t.Fatal(err)
	}

	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		t.Fatal(err)
	}
	if commit.Message != msg {
		t.Errorf("commit message: got %q, want %q", commit.Message, msg)
	}
}
