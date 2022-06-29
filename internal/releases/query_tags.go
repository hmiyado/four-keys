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
	Tag  string
	Date time.Time
}

func (r *Release) String() string {
	return fmt.Sprintf("(Tag=%v,Date=%v)", r.Tag, r.Date)
}

func QueryReleases(repository *git.Repository) []*Release {
	tags := QueryTags(repository)
	releases := make([]*Release, len(tags))
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		commit, _ := repository.CommitObject(tag.Hash())
		releases[i] = &Release{Tag: tag.Name().Short(), Date: commit.Author.When}
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
