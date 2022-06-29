package releases

import (
	"os"
	"sort"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

var r, emptyRepository *git.Repository

func TestMain(m *testing.M) {
	r, _ = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-git",
	})
	emptyRepository, _ = git.Init(memory.NewStorage(), nil)

	code := m.Run()

	os.Exit(code)
}

func TestQueryReleasesShouldHaveSameCountToTags(t *testing.T) {
	releases := QueryReleases(r)
	expectedReleasesNum := len(QueryTags(r))

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
	releases := QueryReleases(r)

	if sort.SliceIsSorted(releases, func(i, j int) bool { return releases[i].Date.Before(releases[j].Date) }) {
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
	releases := QueryReleases(emptyRepository)

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

func TestQueryTagsShouldHaveTags(t *testing.T) {
	tags := QueryTags(r)
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
