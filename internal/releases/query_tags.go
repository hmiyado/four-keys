package releases

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Release struct {
	Tag  string    `json:"tag"`
	Date time.Time `json:"date"`
	// LeadTimeForChanges indicates
	LeadTimeForChanges float64 `json:"leadTimeForChanges"`
}

type Option struct {
	// inclucive
	Since time.Time `json:"since"`
	// inclucive
	Until time.Time `json:"until"`
}

func (r *Release) String() string {
	return fmt.Sprintf("(Tag=%v,Date=%v)", r.Tag, r.Date)
}

func (r *Release) Equal(another *Release) bool {
	return r.Tag == another.Tag && r.Date.Equal(another.Date)
}

func (o *Option) isInTimeRange(time time.Time) bool {
	if o == nil {
		return true
	}
	return time.After(o.Since) && time.Before(o.Until)
}

// QueryReleases returns Releases sorted by date (first item is the oldest and last item is the newest)
func QueryReleases(repository *git.Repository, option *Option) []*Release {
	type ReleaseSource struct {
		tag    *plumbing.Reference
		commit *object.Commit
	}
	tags := QueryTags(repository)
	sources := make([]ReleaseSource, 0)
	for _, tag := range tags {
		commit, err := repository.CommitObject(tag.Hash())
		if err != nil {
			continue
		}
		sources = append(sources, ReleaseSource{tag: tag, commit: commit})
	}
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].commit.Committer.When.Before(sources[j].commit.Committer.When)
	})

	releases := make([]*Release, 0)
	for i, source := range sources {
		if option.isInTimeRange(source.commit.Committer.When) {
			var preCommit *object.Commit
			if i == 0 {
				preCommit = nil
			} else {
				preCommit = sources[i-1].commit
			}
			leadTimeForChanges := GetLeadTimeForChanges(repository, preCommit, source.commit)
			if leadTimeForChanges == nil {
				zero := time.Duration(0)
				leadTimeForChanges = &zero
			}
			releases = append(releases, &Release{
				Tag:                source.tag.Name().Short(),
				Date:               source.commit.Committer.When,
				LeadTimeForChanges: leadTimeForChanges.Hours(),
			})
		}
	}
	return releases
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
