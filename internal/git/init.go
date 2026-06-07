package git

import (
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// InitGit initialises a new git repo in dest and creates a single commit.
func InitGit(dest, message, authorName, authorEmail string, when time.Time) error {
	if dest == "" {
		return fmt.Errorf("dest must not be empty")
	}
	if message == "" {
		return fmt.Errorf("message must not be empty")
	}
	if authorName == "" {
		return fmt.Errorf("author name must not be empty")
	}
	if authorEmail == "" {
		return fmt.Errorf("author email must not be empty")
	}
	if when.IsZero() {
		when = time.Now()
	}

	repo, err := gogit.PlainInit(dest, false)
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
	_, err = w.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  when,
		},
	})
	if err != nil {
		return fmt.Errorf("create initial commit in %s: %w", dest, err)
	}
	return nil
}
