package releases

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Release struct {
	Tag  string    `json:"tag"`
	Date time.Time `json:"date"`
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

func QueryReleases(repository *git.Repository, option *Option) []*Release {
	tags := QueryTags(repository)
	releases := make([]*Release, 0)
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		commit, _ := repository.CommitObject(tag.Hash())
		commitDate := commit.Author.When
		if option.isInTimeRange(commitDate) {
			releases = append(releases, &Release{Tag: tag.Name().Short(), Date: commitDate})
		}
	}
	sort.Slice(releases, func(i, j int) bool { return releases[i].Date.Before(releases[j].Date) })
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
