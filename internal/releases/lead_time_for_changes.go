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
	if newerCommit == nil {
		return nil
	}
	if olderCommit != nil && olderCommit.Hash == newerCommit.Hash {
		zero := time.Duration(0)
		return &zero
	}
	var ancestor, child *object.Commit
	child = newerCommit
	if olderCommit == nil {
		iter, err := repository.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
		if err != nil {
			return nil
		}
		var initialCommit *object.Commit
		iter.ForEach(func(c *object.Commit) error { initialCommit = c; return nil })
		ancestor = initialCommit
	} else {
		ancestor = olderCommit
	}
	if isAncestor, _ := ancestor.IsAncestor(child); !isAncestor {
		return nil
	}
	afterAncestor := ancestor.Committer.When.AddDate(0, 0, 1)
	iter, err := repository.Log(&git.LogOptions{
		From:  child.Hash,
		Since: &afterAncestor,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil
	}
	var lastCommit *object.Commit
	iter.ForEach(func(c *object.Commit) error {
		lastCommit = c
		return nil
	})
	duration := child.Committer.When.Sub(lastCommit.Committer.When)
	return &duration
}
