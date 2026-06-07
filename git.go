package solarcommons

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// InitGit initialises a new git repo in dest and creates a single commit.
func InitGit(dest, message string) error {
	if dest == "" {
		return fmt.Errorf("dest must not be empty")
	}
	if message == "" {
		return fmt.Errorf("message must not be empty")
	}

	repo, err := git.PlainInit(dest, false)
	if err != nil {
		return fmt.Errorf("init git repository in %s: %w", dest, err)
	}
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("open git worktree in %s: %w", dest, err)
	}
	if err := w.AddGlob("."); err != nil {
		return fmt.Errorf("stage files in %s: %w", dest, err)
	}
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "sun",
			Email: "sun@solar",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("create initial commit in %s: %w", dest, err)
	}
	return nil
}
