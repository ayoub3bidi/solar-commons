package solarcommons

import (
	gitops "github.com/ayoub3bidi/solar-commons/internal/git"
)

// CloneTemplate shallow-clones repo at the given tag into dest.
//
// Template manifests should pin stable semver tags for end-user scaffolds.
func CloneTemplate(repo, ref, dest string) error {
	return gitops.CloneTemplate(repo, ref, dest)
}

// StripGit removes the .git directory from dest.
func StripGit(dest string) error {
	return gitops.StripGit(dest)
}
