package core

import "testing"

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

func TestQueryTagsShouldRecognizeTagObject(t *testing.T) {
	tags := QueryTags(repositoryCli)
	expectedTagsCount := 77

	if len(tags) < expectedTagsCount {
		// if QueryTags can't recognize tag object,
		// tags count becomes much smaller than expectedTagsCount
		t.Errorf("cli repository should have over %v tags but %v", expectedTagsCount, len(tags))
	}
}
