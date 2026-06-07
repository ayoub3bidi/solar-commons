package solarcommons

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneTemplate shallow-clones repo at the given tag into dest.
//
// Template manifests should pin stable semver tags for end-user scaffolds.
func CloneTemplate(repo, ref, dest string) error {
	if repo == "" {
		return fmt.Errorf("repo must not be empty")
	}
	if ref == "" {
		return fmt.Errorf("ref must not be empty")
	}
	if dest == "" {
		return fmt.Errorf("dest must not be empty")
	}

	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:           repo,
		ReferenceName: plumbing.NewTagReferenceName(ref),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return fmt.Errorf("clone %s@%s into %s: %w", repo, ref, dest, err)
	}
	return nil
}

// StripGit removes the .git directory from dest.
func StripGit(dest string) error {
	if dest == "" {
		return fmt.Errorf("dest must not be empty")
	}
	if err := os.RemoveAll(filepath.Join(dest, ".git")); err != nil {
		return fmt.Errorf("strip git directory from %s: %w", dest, err)
	}
	return nil
}
