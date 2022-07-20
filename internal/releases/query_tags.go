package releases

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Release struct {
	Tag                string        `json:"tag"`
	Date               time.Time     `json:"date"`
	LeadTimeForChanges time.Duration `json:"leadTimeForChanges"`
	Result             ReleaseResult `json:"result"`
}

type Option struct {
	// inclucive
	Since time.Time `json:"since"`
	// inclucive
	Until          time.Time      `json:"until"`
	IgnorePattern  *regexp.Regexp `json:"ignorePattern"`
	StartTimerFunc func(string)   `json:"-"`
	StopTimerFunc  func(string)   `json:"-"`
}

func (r *Release) String() string {
	return fmt.Sprintf("(Tag=%v, Date=%v, LeadTimeForChamges=%v, Result=%v)", r.Tag, r.Date, r.LeadTimeForChanges, r.Result.String())
}

func (r *Release) Equal(another *Release) bool {
	return r.Tag == another.Tag && r.Date.Equal(another.Date) && r.Result.Equal(another.Result)
}

func (o *Option) isInTimeRange(time time.Time) bool {
	if o == nil {
		return true
	}
	return time.After(o.Since) && time.Before(o.Until)
}

func (o *Option) shouldIgnore(name string) bool {
	if o == nil || o.IgnorePattern == nil {
		return false
	}
	return o.IgnorePattern.MatchString(name)
}

func (o *Option) StartTimer(key string) {
	if o != nil && o.StartTimerFunc != nil {
		o.StartTimerFunc(key)
	}
}

func (o *Option) StopTimer(key string) {
	if o != nil && o.StopTimerFunc != nil {
		o.StopTimerFunc(key)
	}
}

// QueryReleases returns Releases sorted by date (first item is the oldest and last item is the newest)
func QueryReleases(repository *git.Repository, option *Option) []*Release {
	option.StartTimer("QueryReleases")
	defer option.StopTimer("QueryReleases")
	option.StartTimer("QueryTags")
	tags := QueryTags(repository)
	option.StopTimer("QueryTags")
	sources := getReleaseSourcesFromTags(repository, tags)

	releases := make([]*Release, 0)
	for i, source := range sources {
		if option.shouldIgnore(source.tag.Name().Short()) {
			continue
		}

		if !option.isInTimeRange(source.commit.Committer.When) {
			continue
		}

		timerKeyGetLeadTimeForChanges := fmt.Sprintf("source[%v]GetLeadTimeForChanges", i)
		option.StartTimer(timerKeyGetLeadTimeForChanges)
		var preCommit *object.Commit
		if i == len(sources)-1 {
			preCommit = nil
		} else {
			preCommit = sources[i+1].commit
		}
		leadTimeForChanges := GetLeadTimeForChanges(repository, preCommit, source.commit)
		if leadTimeForChanges == nil {
			zero := time.Duration(0)
			leadTimeForChanges = &zero
		}
		option.StopTimer(timerKeyGetLeadTimeForChanges)

		timerKeyGetReleaseResult := fmt.Sprintf("source[%v]ReleaseResult", i)
		option.StartTimer(timerKeyGetReleaseResult)
		releases = append(releases, &Release{
			Tag:                source.tag.Name().Short(),
			Date:               source.commit.Committer.When,
			LeadTimeForChanges: *leadTimeForChanges,
			Result:             getReleaseResult(repository, sources, i),
		})
		option.StopTimer(timerKeyGetReleaseResult)
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
