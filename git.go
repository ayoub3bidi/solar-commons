package solarcommons

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// InitGit initialises a new git repo in dest and creates a single commit.
func InitGit(dest, message string) error {
	repo, err := git.PlainInit(dest, false)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	if err := w.AddGlob("."); err != nil {
		return err
	}
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "sun",
			Email: "sun@solar",
			When:  time.Now(),
		},
	})
	return err
}
