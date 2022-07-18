package releases

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type ReleaseResult struct {
	IsSuccess     bool
	TimeToRestore *time.Duration
}

func (r ReleaseResult) Equal(another ReleaseResult) bool {
	if r.TimeToRestore == nil && another.TimeToRestore == nil {
		return r.IsSuccess == another.IsSuccess
	}
	if r.TimeToRestore != nil && another.TimeToRestore != nil {
		return r.IsSuccess == another.IsSuccess && *(r.TimeToRestore) == *(another.TimeToRestore)
	}
	return false
}

func (r ReleaseResult) String() string {
	return fmt.Sprintf("IsSuccess=%v, TimeToRestore=%v", r.IsSuccess, r.TimeToRestore)
}

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

// getReleaseResult returns ReleaseResult for targetIndex with repository and sources.
// It is assumed that sources is sorted by descending order of date.
func getReleaseResult(repository *git.Repository, sources []ReleaseSource, targetIndex int) ReleaseResult {
	if targetIndex == 0 {
		// it is considered that the newest release is success
		return ReleaseResult{
			IsSuccess:     true,
			TimeToRestore: nil,
		}
	}
	source := sources[targetIndex]
	isRestored := isRestoredRelease(repository, source.commit, sources[targetIndex-1].commit)
	isSuccess := !isRestored

	var preReleaseCommit *object.Commit
	if targetIndex < len(sources)-1 {
		preReleaseCommit = sources[targetIndex+1].commit
	}

	timeToRestore := time.Duration(0)
	if isSuccess && isRestoredRelease(repository, preReleaseCommit, source.commit) {
		newerCommitIndex := targetIndex
		for isRestoredRelease(repository, sources[newerCommitIndex+1].commit, sources[newerCommitIndex].commit) {
			timeToRestore += sources[newerCommitIndex].commit.Committer.When.Sub(sources[newerCommitIndex+1].commit.Committer.When)
			newerCommitIndex += 1
			if newerCommitIndex >= len(sources)-1 {
				break
			}
		}
	}

	var resultTimeToRestore *time.Duration
	if timeToRestore != 0 {
		resultTimeToRestore = &timeToRestore
	}

	return ReleaseResult{
		IsSuccess:     isSuccess,
		TimeToRestore: resultTimeToRestore,
	}
}
