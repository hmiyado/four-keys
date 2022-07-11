package releases

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
)

func TestGetLeadTimeForChangesShouldReturnLeadTime(t *testing.T) {
	// https://github.com/go-git/go-git/releases/tag/v5.0.0
	commit5_0_0, err1 := repository.CommitObject(plumbing.NewHash("9d0f15c4fa712cdacfa3887e9baac918f093fbf6"))
	if err1 != nil {
		t.Error(err1.Error())
	}
	// https://github.com/go-git/go-git/releases/tag/v5.1.0
	commit5_1_0, err2 := repository.CommitObject(plumbing.NewHash("8019144b6534ff58ad234a355e5b143f1c99b45e"))
	if err2 != nil {
		t.Error(err1.Error())
	}

	duration := GetLeadTimeForChanges(repository, commit5_0_0, commit5_1_0)
	expected, _ := time.ParseDuration("1677h6m36s")
	if isNearBy(duration.Minutes(), expected.Minutes(), 0.01) {
		return
	}
	t.Error("should be x0.99-1.01 of ", expected, " but ", duration)
}

func TestGetLeadTimeForChangesShouldReturnSameTimeForSwitchingCommits(t *testing.T) {
	// https://github.com/go-git/go-git/releases/tag/v5.0.0
	commit5_0_0, err1 := repository.CommitObject(plumbing.NewHash("9d0f15c4fa712cdacfa3887e9baac918f093fbf6"))
	if err1 != nil {
		t.Error(err1.Error())
	}
	// https://github.com/go-git/go-git/releases/tag/v5.1.0
	commit5_1_0, err2 := repository.CommitObject(plumbing.NewHash("8019144b6534ff58ad234a355e5b143f1c99b45e"))
	if err2 != nil {
		t.Error(err1.Error())
	}

	duration1 := GetLeadTimeForChanges(repository, commit5_0_0, commit5_1_0)
	duration2 := GetLeadTimeForChanges(repository, commit5_1_0, commit5_0_0)
	if duration1.Minutes() == duration2.Minutes() {
		return
	}
	t.Error(duration1, "should be equal to ", duration2)
}

func TestGetLeadTimeForChangesShouldReturn0ForSameCommits(t *testing.T) {
	// https://github.com/go-git/go-git/releases/tag/v5.0.0
	commit5_0_0, err1 := repository.CommitObject(plumbing.NewHash("9d0f15c4fa712cdacfa3887e9baac918f093fbf6"))
	if err1 != nil {
		t.Error(err1.Error())
	}

	duration := GetLeadTimeForChanges(repository, commit5_0_0, commit5_0_0)
	if duration.Minutes() == 0 {
		return
	}
	t.Error(duration, "should be equal to ", 0)
}

// isNearBy checks actual is in range of [expected*(1-epsilon), expected*(1+epsiolon)]
func isNearBy(actual float64, expected float64, epsilon float64) bool {
	return actual >= expected*(1-epsilon) && actual <= expected*(1+epsilon)
}
