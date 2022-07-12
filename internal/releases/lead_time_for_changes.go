package releases

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetLeadTimeForChanges returns duration after older commit is created (exlusive) and until newer commit is created (inclusive).
// if commit1 xor commit2 is nil, it is assumed that older commit is the initial commit of this repository.
// if commit1 and commit2 is nil, return duration of 0.
func GetLeadTimeForChanges(repository *git.Repository, commit1 *object.Commit, commit2 *object.Commit) *time.Duration {
	if commit1 == nil && commit2 == nil {
		zero := time.Duration(0)
		return &zero
	}
	if commit1 == commit2 {
		zero := time.Duration(0)
		return &zero
	}
	var ancestor, child *object.Commit
	if commit1 == nil || commit2 == nil {
		iter, err := repository.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
		if err != nil {
			return nil
		}
		var initialCommit *object.Commit
		iter.ForEach(func(c *object.Commit) error { initialCommit = c; return nil })
		ancestor = initialCommit
		if commit1 == nil {
			child = commit2
		} else {
			child = commit1
		}
	} else if isAncestor, _ := commit1.IsAncestor(commit2); isAncestor {
		ancestor = commit1
		child = commit2
	} else if isAncestor, _ := commit2.IsAncestor(commit1); isAncestor {
		ancestor = commit2
		child = commit1
	} else {
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
