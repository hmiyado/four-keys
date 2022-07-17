package releases

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// isResoredRelease returns true if newerCommit restores olderCommit.
// If a commit containing 'hotfix' exists between newerCommit and olderCommit, it is considered that newerCommit restores olderCommit .
// It returns nil if newer commit is null.
// It returns nil if older commit is not ancestor of newer.
// If older commit is null, it is assumed that older commit is the initial commit of this repository.
func isRestoredRelease(repository *git.Repository, olderCommit *object.Commit, newerCommit *object.Commit) bool {
	isRestored := false
	traverseCommits(repository, olderCommit, newerCommit, func(c *object.Commit) error {
		if strings.Contains(c.Message, "hotfix") {
			isRestored = true
			return storer.ErrStop
		}
		return nil
	})
	return isRestored
}
