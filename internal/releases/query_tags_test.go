package releases

import (
	"os"
	"regexp"
	"sort"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

var repository, emptyRepository *git.Repository

func TestMain(m *testing.M) {
	repository, _ = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-git",
	})
	emptyRepository, _ = git.Init(memory.NewStorage(), nil)

	code := m.Run()

	os.Exit(code)
}

func TestQueryReleasesShouldHaveSameCountToTags(t *testing.T) {
	releases := QueryReleases(repository, nil)
	expectedReleasesNum := len(QueryTags(repository))

	if len(releases) != expectedReleasesNum {
		for i := 0; i < len(releases); i++ {
			if releases[i] == nil {
				t.Logf("releases[%v] = nil", i)
			} else {
				t.Logf("releases[%d] = %s", i, releases[i].String())
			}
		}
		t.Errorf("num of tags should be %d but %d", expectedReleasesNum, len(releases))
	}
}

func TestQueryReleasesShouldBeSortedByDate(t *testing.T) {
	releases := QueryReleases(repository, nil)

	if sort.SliceIsSorted(releases, func(i, j int) bool { return releases[i].Date.After(releases[j].Date) }) {
		return
	}
	for i := 0; i < len(releases); i++ {
		if releases[i] == nil {
			t.Logf("releases[%v] = nil", i)
		} else {
			t.Logf("releases[%d] = %s", i, releases[i].String())
		}
	}
	t.Errorf("releases are not sorted")
}

func TestQueryReleasesShouldReturnEmptyForEmptyRepository(t *testing.T) {
	releases := QueryReleases(emptyRepository, nil)

	if len(releases) == 0 {
		return
	}
	for i := 0; i < len(releases); i++ {
		if releases[i] == nil {
			t.Logf("releases[%v] = nil", i)
		} else {
			t.Logf("releases[%d] = %s", i, releases[i].String())
		}
	}
	t.Errorf("releases are not sorted")
}

func TestQueryReleasesShouldReturnReleasesWithSpecifiedTimeRange(t *testing.T) {
	releases := QueryReleases(repository, &Option{
		Since: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Until: time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC),
	})
	tag5_0_0 := &Release{Tag: "v5.0.0", Date: time.Date(2020, 3, 15, 21, 18, 32, 0, time.FixedZone("+0100", 1*60*60)), Result: ReleaseResult{IsSuccess: true}}
	tag5_1_0 := &Release{Tag: "v5.1.0", Date: time.Date(2020, 5, 24, 19, 25, 8, 0, time.FixedZone("+0200", 2*60*60)), Result: ReleaseResult{IsSuccess: true}}
	tag5_2_0 := &Release{Tag: "v5.2.0", Date: time.Date(2020, 10, 9, 11, 49, 30, 0, time.FixedZone("+0200", 2*60*60)), Result: ReleaseResult{IsSuccess: true}}
	expectedTags := []*Release{tag5_2_0, tag5_1_0, tag5_0_0}

	if len(releases) != len(expectedTags) {
		t.Errorf("releases does not have expected tag num. expected: %v. actual: %v", len(expectedTags), len(releases))
		return
	}

	unmatchedRelease := make([]int, 0)
	for i, actual := range releases {
		expected := expectedTags[i]
		if actual.Equal(expected) {
			continue
		}
		unmatchedRelease = append(unmatchedRelease, i)
	}

	if len(unmatchedRelease) == 0 {
		return
	}

	for i := range unmatchedRelease {
		actual := releases[i]
		expected := expectedTags[i]
		t.Logf("releases[%d] = %s. expected: %v", i, actual, expected)
	}
	t.Errorf("releases does not have specified")
}

func TestQueryReleasesShouldHaveReleaseResult(t *testing.T) {
	releases := QueryReleases(repository, &Option{
		Since: time.Date(2015, 12, 20, 0, 0, 0, 0, time.UTC),
		Until: time.Date(2016, 1, 11, 23, 59, 59, 999, time.UTC),
	})
	tag2_1_2 := &Release{
		Tag:                "v2.1.2",
		Date:               time.Date(2016, 1, 11, 12, 9, 15, 0, time.FixedZone("+0100", 1*60*60)),
		LeadTimeForChanges: parseDurationOrZero("25m24s"),
		Result: ReleaseResult{
			IsSuccess:            true,
			TimeToRestoreService: parseDurationOrNil("67h7m39s"),
		},
	}
	tag2_1_1 := &Release{
		Tag:                "v2.1.1",
		Date:               time.Date(2016, 1, 8, 17, 1, 36, 0, time.FixedZone("+0100", 1*60*60)),
		LeadTimeForChanges: parseDurationOrZero("12m26s"),
		Result: ReleaseResult{
			IsSuccess:            false,
			TimeToRestoreService: nil,
		},
	}
	tag2_1_0 := &Release{
		Tag:                "v2.1.0",
		Date:               time.Date(2015, 12, 23, 9, 48, 11, 0, time.FixedZone("+0100", 1*60*60)),
		LeadTimeForChanges: parseDurationOrZero(""),
		Result: ReleaseResult{
			IsSuccess:            true,
			TimeToRestoreService: nil,
		},
	}
	expectedTags := []*Release{tag2_1_2, tag2_1_1, tag2_1_0}

	if len(releases) != len(expectedTags) {
		t.Errorf("releases does not have expected tag num. expected: %v. actual: %v", len(expectedTags), len(releases))
		return
	}

	unmatchedRelease := make([]int, 0)
	for i, actual := range releases {
		expected := expectedTags[i]
		if actual.Equal(expected) {
			continue
		}
		unmatchedRelease = append(unmatchedRelease, i)
	}

	if len(unmatchedRelease) == 0 {
		return
	}

	for i := range unmatchedRelease {
		actual := releases[i]
		expected := expectedTags[i]
		t.Logf("releases[%d] = %s. expected: %v", i, actual, expected)
	}
	t.Errorf("releases does not have specified")
}

func TestQueryReleasesShouldReturnReleasesWithIgnorePattern(t *testing.T) {
	pattern, _ := regexp.Compile(`v5\.1\.0|v5\.2\.0`)
	releases := QueryReleases(repository, &Option{
		Since:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:         time.Date(2020, 12, 31, 23, 59, 59, 999, time.UTC),
		IgnorePattern: pattern,
	})
	tag5_0_0 := &Release{Tag: "v5.0.0", Date: time.Date(2020, 3, 15, 21, 18, 32, 0, time.FixedZone("+0100", 1*60*60))}
	expectedTags := []*Release{tag5_0_0}

	if len(releases) != len(expectedTags) {
		t.Errorf("releases does not have expected tag num. expected: %v. actual: %v", len(expectedTags), len(releases))
		return
	}

	if releases[0].Tag == tag5_0_0.Tag {
		return
	}

	t.Logf("expected: %v, actual: %v", tag5_0_0, releases[0])
	t.Errorf("releases does not have specified")
}

func TestQueryTagsShouldHaveTags(t *testing.T) {
	tags := QueryTags(repository)
	expectedTagNum := 60

	if len(tags) != expectedTagNum {
		for i := 0; i < len(tags); i++ {
			if tags[i] == nil {
				t.Logf("tags[%d] = nil", i)
			} else {
				t.Logf("tags[%d] = %s", i, tags[i].String())
			}
		}
		t.Errorf("num of tags should be %d but %d", expectedTagNum, len(tags))
	}
}

func TestQueryTagsShouldReturnEmptyForEmptyRepository(t *testing.T) {
	tags := QueryTags(emptyRepository)
	expectedTagNum := 0

	if len(tags) != expectedTagNum {
		for i := 0; i < len(tags); i++ {
			if tags[i] == nil {
				t.Logf("tags[%d] = nil", i)
			} else {
				t.Logf("tags[%d] = %s", i, tags[i].String())
			}
		}
		t.Errorf("num of tags should be %d but %d", expectedTagNum, len(tags))
	}
}

func parseDurationOrZero(str string) time.Duration {
	d, err := time.ParseDuration(str)
	if err != nil {
		return time.Duration(0)
	}
	return d
}

func parseDurationOrNil(str string) *time.Duration {
	d, err := time.ParseDuration(str)
	if err != nil {
		return nil
	}
	return &d
}
