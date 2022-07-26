package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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
	Until            time.Time      `json:"until"`
	IgnorePattern    *regexp.Regexp `json:"ignorePattern"`
	FixCommitPattern *regexp.Regexp
	StartTimerFunc   func(string) `json:"-"`
	StopTimerFunc    func(string) `json:"-"`
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

func (o *Option) isFixedCommit(commitMessage string) bool {
	if o == nil || o.FixCommitPattern == nil {
		// commitMessage with "hotfix" is regarded as fixed commit by default
		return strings.Contains(commitMessage, "hotfix")
	}
	return o.FixCommitPattern.MatchString(commitMessage)
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
	nextSuccessSourceIndex := -1
	nextSuccessReleaseIndex := -1
	isRestored := false
	for i, source := range sources {
		if option.shouldIgnore(source.tag.Name().Short()) {
			continue
		}

		if !option.isInTimeRange(source.commit.Committer.When) {
			continue
		}

		timerKeyReleaseMetrics := fmt.Sprintf("source[%v]GetReleaseMetrics", i)
		option.StartTimer(timerKeyReleaseMetrics)

		isSuccess := !isRestored
		if nextSuccessSourceIndex > 0 && isSuccess && i != nextSuccessSourceIndex+1 {
			timeToRestore := sources[nextSuccessSourceIndex].commit.Committer.When.Sub(sources[i-1].commit.Committer.When)
			releases[nextSuccessReleaseIndex].Result.TimeToRestore = &timeToRestore
		}
		if isSuccess {
			nextSuccessSourceIndex = i
			nextSuccessReleaseIndex = len(releases)
		}

		var lastCommit *object.Commit
		var preReleaseCommit *object.Commit
		if i < len(sources)-1 {
			preReleaseCommit = sources[i+1].commit
		}
		restoresPreRelease := false
		err := traverseCommits(repository, preReleaseCommit, source.commit, func(c *object.Commit) error {
			if option.isFixedCommit(c.Message) {
				restoresPreRelease = true
			}
			lastCommit = c
			return nil
		})
		isRestored = restoresPreRelease
		leadTimeForChanges := time.Duration(0)
		if err == nil && lastCommit != nil {
			leadTimeForChanges = source.commit.Committer.When.Sub(lastCommit.Committer.When)
		}
		option.StopTimer(timerKeyReleaseMetrics)

		releases = append(releases, &Release{
			Tag:                source.tag.Name().Short(),
			Date:               source.commit.Committer.When,
			LeadTimeForChanges: leadTimeForChanges,
			Result: ReleaseResult{
				IsSuccess: isSuccess,
			},
		})
	}
	return releases
}
