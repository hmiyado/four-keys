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
