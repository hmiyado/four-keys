package releases

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestQueryTags(t *testing.T) {
	r, _ := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-git",
	})

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
