package core

import (
	"errors"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

type ReleaseSource struct {
	tag    *plumbing.Reference
	commit *object.Commit
}

func getReleaseSourcesFromTags(repository *git.Repository, tags []*plumbing.Reference) []ReleaseSource {
	sources := make([]ReleaseSource, 0)
	for _, tag := range tags {
		commit, err := repository.CommitObject(tag.Hash())
		if err != nil {
			continue
		}
		sources = append(sources, ReleaseSource{tag: tag, commit: commit})
	}
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].commit.Committer.When.After(sources[j].commit.Committer.When)
	})
	return sources
}

func QueryTags(repository *git.Repository) []*plumbing.Reference {
	itr, _ := repository.Tags()
	tags := make([]*plumbing.Reference, 0)

	itr.ForEach(func(ref *plumbing.Reference) error {
		// refs/tags/xxx
		// lightweight tag
		if strings.Split(ref.Name().String(), "/")[1] == "tags" {
			tags = append(tags, ref)
			return nil
		}

		_, err := repository.TagObject(ref.Hash())
		if err != nil {
			tags = append(tags, ref)
		}
		return nil
	})

	return tags
}
