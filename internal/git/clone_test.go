package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ayoub3bidi/solar-commons/internal/testutil"
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
	source := testutil.NewTaggedRepo(t)
	dest := t.TempDir()

	if err := CloneTemplate(source, "v0.1.0", dest); err != nil {
		t.Fatalf("CloneTemplate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); os.IsNotExist(err) {
		t.Fatal(".git should exist after clone")
	}
	if err := StripGit(dest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".git")); !os.IsNotExist(err) {
		t.Fatal(".git should be gone after StripGit")
	}
}
