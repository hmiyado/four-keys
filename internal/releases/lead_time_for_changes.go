package releases

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetLeadTimeForChanges(repository *git.Repository, commit1 *object.Commit, commit2 *object.Commit) *time.Duration {
	var ancestor, child *object.Commit
	if isAncestor, _ := commit1.IsAncestor(commit2); isAncestor {
		ancestor = commit1
		child = commit2
	} else if isAncestor, _ := commit2.IsAncestor(commit1); isAncestor {
		ancestor = commit2
		child = commit1
	} else {
		return nil
	}
	iter, err := repository.Log(&git.LogOptions{
		From:  child.Hash,
		Since: &ancestor.Committer.When,
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
