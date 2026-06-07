package solarcommons

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneTemplate shallow-clones repo at the given ref into dest.
func CloneTemplate(repo, ref, dest string) error {
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:           repo,
		ReferenceName: plumbing.NewTagReferenceName(ref),
		SingleBranch:  true,
		Depth:         1,
	})
	return err
}

// StripGit removes the .git directory from dest.
func StripGit(dest string) error {
	return os.RemoveAll(filepath.Join(dest, ".git"))
}
