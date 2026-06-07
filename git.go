package solarcommons

import (
	"time"

	gitops "github.com/ayoub3bidi/solar-commons/internal/git"
)

// InitGit initialises a new git repo in dest and creates a single commit.
func InitGit(dest, message string) error {
	return gitops.InitGit(dest, message, "sun", "sun@solar", time.Now())
}
