package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
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
	isRestored         bool          `json:-`
}

type Option struct {
	// inclucive
	Since time.Time `json:"since"`
	// inclucive
	Until             time.Time      `json:"until"`
	IgnorePattern     *regexp.Regexp `json:"-"`
	FixCommitPattern  *regexp.Regexp `json:"-"`
	IsLocalRepository bool           `json:"-"`
	StartTimerFunc    func(string)   `json:"-"`
	StopTimerFunc     func(string)   `json:"-"`
	DebuglnFunc       func(...any)   `json:"-"`
}

func (r *Release) String() string {
	return fmt.Sprintf("(Tag=%v, Date=%v, LeadTimeForChamges=%v, Result=%v)", r.Tag, r.Date, r.LeadTimeForChanges.Nanoseconds(), r.Result.String())
}

func (r *Release) Equal(another *Release) bool {
	return r.Tag == another.Tag &&
		r.Date.Equal(another.Date) &&
		r.LeadTimeForChanges.Nanoseconds() == another.LeadTimeForChanges.Nanoseconds() &&
		r.Result.Equal(another.Result)
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

func (o *Option) Debugln(a ...any) {
	if o != nil && o.DebuglnFunc != nil {
		o.DebuglnFunc(a...)
	}
}

// ignoreReleases filters sources by option.IgnorePattern.
// This doesn't filter any sources and returns same sources if no ignore option.
func ignoreReleases(sources []ReleaseSource, option *Option) []ReleaseSource {
	if option == nil || option.IgnorePattern == nil {
		return sources
	}
	filteredSources := make([]ReleaseSource, 0)
	for _, source := range sources {
		if option.shouldIgnore(source.tag.Name().Short()) {
			option.Debugln(source.tag.Name().Short(), " is ignored")
			continue
		}
		filteredSources = append(filteredSources, source)
	}
	return filteredSources
}

func createReleasesBySources(sources []ReleaseSource, option *Option, repository *git.Repository) []*Release {
	releases := make([]*Release, 0)

	for i, source := range sources {
		if !option.isInTimeRange(source.commit.Committer.When) {
			option.Debugln("source[", i, "](", source.tag.Name().Short(), ") is skipped for outof time range")
			continue
		}

		timerEachReleases := fmt.Sprintf("source[%v](%v)GetEachReleases", i, source.tag.Name().Short())
		option.StartTimer(timerEachReleases)

		leadTimeForChanges := time.Duration(0)
		isRestored := false
		if option != nil && option.IsLocalRepository {
			isRestored, leadTimeForChanges = getIsRestoredAndLeadTimeForChangesByLocalGit(sources, i, option)
		} else {
			isRestored, leadTimeForChanges = getIsRestoredAndLeadTimeForChangesByGoGit(sources, i, option, repository)
		}
		option.StopTimer(timerEachReleases)

		releases = append(releases, &Release{
			Tag:                source.tag.Name().Short(),
			Date:               source.commit.Committer.When,
			LeadTimeForChanges: leadTimeForChanges,
			Result: ReleaseResult{
				IsSuccess: false,
			},
			isRestored: isRestored,
		})

	}

	return releases
}

func setReleaseResultForEachRelease(releases []*Release, option *Option) {
	var nextSuccessRelease *Release
	for i, release := range releases {
		if !option.isInTimeRange(release.Date) {
			option.Debugln("release[", i, "](", release.Tag, ") is skipped for outof time range")
			continue
		}

		timerKeyReleaseMetrics := fmt.Sprintf("release[%v](%v)GetReleaseMetrics", i, release.Tag)
		option.StartTimer(timerKeyReleaseMetrics)

		isSuccess := true
		var nextRelease *Release
		if i > 0 {
			nextRelease = releases[i-1]
			isSuccess = !nextRelease.isRestored
		}
		release.Result.IsSuccess = isSuccess
		if isSuccess {
			if nextRelease != nil && !nextRelease.Result.IsSuccess {
				timeToRestore := nextSuccessRelease.Date.Sub(releases[i-1].Date)
				nextSuccessRelease.Result.TimeToRestore = &timeToRestore
			}
			nextSuccessRelease = release
		}

		option.StopTimer(timerKeyReleaseMetrics)
	}
}

// QueryReleases returns Releases sorted by date (first item is the oldest and last item is the newest)
func QueryReleases(repository *git.Repository, option *Option) []*Release {
	option.StartTimer("QueryReleases")
	defer option.StopTimer("QueryReleases")
	option.StartTimer("QueryTags")
	tags := QueryTags(repository)
	option.StopTimer("QueryTags")
	option.Debugln("Tags count:", len(tags))
	sources := getReleaseSourcesFromTags(repository, tags)
	sources = ignoreReleases(sources, option)
	option.Debugln("Sources count:", len(sources))

	releases := createReleasesBySources(sources, option, repository)
	setReleaseResultForEachRelease(releases, option)
	return releases
}

// getIsRestoredAndLeadTimeForChangesByLocalGit gets isRestored and leadTimeForChanges by using local git command.
// Local git command is about 10 times faster than go-git.
func getIsRestoredAndLeadTimeForChangesByLocalGit(
	sources []ReleaseSource,
	i int,
	option *Option,
) (isRestored bool, leadTimeForChanges time.Duration) {
	source := sources[i]
	since := "1900-01-01"
	if i < len(sources)-1 {
		preReleaseCommit := sources[i+1].commit
		since = preReleaseCommit.Committer.When.Add(time.Second).Format("2006-01-02T15:04:05")
	}
	restoresPreRelease := false
	output, cmdErr := exec.Command("git", "log",
		"--since", since,
		`--format="%ct %s"`,
		"--date-order",
		source.commit.Hash.String(),
	).Output()
	lastLine := ""
	if cmdErr == nil {
		scanner := bufio.NewScanner(bytes.NewReader(output))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			if option.isFixedCommit(line) {
				restoresPreRelease = true
			}
			lastLine = line
		}
	}
	isRestored = restoresPreRelease
	leadTimeForChanges = time.Duration(0)
	pattern, _ := regexp.Compile(`^[^\d]+(\d+) `)
	if cmdErr == nil && lastLine != "" {
		matches := pattern.FindStringSubmatch(lastLine)
		unixtimeString := matches[1]
		unixtimeInt, err := strconv.ParseInt(unixtimeString, 10, 64)
		if err == nil {
			lastCommitWhen := time.Unix(unixtimeInt, 0)
			leadTimeForChanges = source.commit.Committer.When.Sub(lastCommitWhen)
		}
	}
	return isRestored, leadTimeForChanges
}

// getIsRestoredAndLeadTimeForChangesByGoGit gets isRestored and leadTimeForChanges by using go-git.
// go-git is slow but it can use in-memory repository.
// When repository is specified by url, repository is in-memory so that go-git is used.
func getIsRestoredAndLeadTimeForChangesByGoGit(
	sources []ReleaseSource,
	i int,
	option *Option,
	repository *git.Repository,
) (isRestored bool, leadTimeForChanges time.Duration) {
	source := sources[i]
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
	leadTimeForChanges = time.Duration(0)
	if err == nil && lastCommit != nil {
		leadTimeForChanges = source.commit.Committer.When.Sub(lastCommit.Committer.When)
	}
	return isRestored, leadTimeForChanges
}
