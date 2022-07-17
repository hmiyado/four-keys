package releases

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrNoNewerCommit      = errors.New("no newer commit")
	ErrLogUnavailable     = errors.New("repository log is unavailable")
	ErrIllegalCommitOrder = errors.New("newer commit is older than older commit")
)

// traverseCommits runs traversaler for each commits between olderCommit and newerCommit.
// Order to traverse rely on git.LogOrderCommitterTime.
// If older commit is null, it is assumed that older commit is the initial commit of this repository.
// It returns ErrNoNewerCommit if newer commit is null.
// It returns ErrIllegalCommitOrder if older commit is not ancestor of newer.
// It returns ErrLogUnavailable if git log is unavailable.
func traverseCommits(repository *git.Repository, olderCommit *object.Commit, newerCommit *object.Commit, traversaler func(c *object.Commit) error) error {
	if newerCommit == nil {
		return ErrNoNewerCommit
	}
	var ancestor, child *object.Commit
	child = newerCommit
	if olderCommit == nil {
		iter, err := repository.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
		if err != nil {
			return ErrLogUnavailable
		}
		var initialCommit *object.Commit
		iter.ForEach(func(c *object.Commit) error { initialCommit = c; return nil })
		ancestor = initialCommit
	} else {
		ancestor = olderCommit
	}
	if isAncestor, _ := ancestor.IsAncestor(child); !isAncestor {
		return ErrIllegalCommitOrder
	}
	afterAncestor := ancestor.Committer.When.AddDate(0, 0, 1)
	iter, err := repository.Log(&git.LogOptions{
		From:  child.Hash,
		Since: &afterAncestor,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return ErrLogUnavailable
	}
	return iter.ForEach(traversaler)
}
