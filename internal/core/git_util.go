package core

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var (
	ErrNoNewerCommit  = errors.New("no newer commit")
	ErrLogUnavailable = errors.New("repository log is unavailable")
)

// traverseCommits runs traversaler for each commits between olderCommit and newerCommit.
// Order to traverse rely on git.LogOrderCommitterTime.
// If older commit is null, it is assumed that older commit is the initial commit of this repository.
// It returns ErrNoNewerCommit if newer commit is null.
// It returns ErrLogUnavailable if git log is unavailable.
func traverseCommits(repository *git.Repository, olderCommit *object.Commit, newerCommit *object.Commit, traversaler func(c *object.Commit) error) error {
	if newerCommit == nil {
		return ErrNoNewerCommit
	}
	var since *time.Time
	if olderCommit != nil {
		tmp := olderCommit.Committer.When.AddDate(0, 0, 1)
		since = &tmp
	}
	iter, err := repository.Log(&git.LogOptions{
		From:  newerCommit.Hash,
		Since: since,
		// Order: git.LogOrderDefault,
		// Order: git.LogOrderDFS,
		// Order: git.LogOrderDFSPost,
		// Order: git.LogOrderBSF,
		// Order: git.LogOrderCommitterTime,
		// Order: git.LogOrderCommitterTime,
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
			tagObject, err := repository.TagObject(tag.Hash())
			if err != nil {
				continue
			}
			commit, err = tagObject.Commit()
			if err != nil {
				continue
			}
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
