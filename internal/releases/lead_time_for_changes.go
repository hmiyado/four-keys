package releases

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetLeadTimeForChanges returns duration after older commit is created (exlusive) and until newer commit is created (inclusive).
// It returns nil if newer commit is null.
// It returns 0 if newer commit is equal to older commit.
// It returns nil if older commit is not ancestor of newer.
// If older commit is null, it is assumed that older commit is the initial commit of this repository.
func GetLeadTimeForChanges(repository *git.Repository, olderCommit *object.Commit, newerCommit *object.Commit) *time.Duration {
	if olderCommit != nil && olderCommit.Hash == newerCommit.Hash {
		zero := time.Duration(0)
		return &zero
	}

	var lastCommit *object.Commit
	err := traverseCommits(repository, olderCommit, newerCommit, func(c *object.Commit) error {
		lastCommit = c
		return nil
	})
	if err != nil {
		return nil
	}
	if lastCommit == nil {
		return nil
	}
	duration := newerCommit.Committer.When.Sub(lastCommit.Committer.When)
	return &duration
}
