package releases

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hmiyado/four-keys/internal/util"
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
	expected, _ := time.ParseDuration("1605h57m40s")
	if util.IsNearBy(duration.Minutes(), expected.Minutes(), 0.01) {
		return
	}
	t.Error("should be x0.99-1.01 of ", expected, " but ", duration)
}

func TestGetLeadTimeForChangesShouldReturnNilWhenNewerCommitIsNotNewer(t *testing.T) {
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

	// commit5_1_0 is newer than commit5_0_0, but argument is swapped newer and older
	duration := GetLeadTimeForChanges(repository, commit5_1_0, commit5_0_0)
	if duration == nil {
		return
	}
	t.Error("should be equal to nil but ", duration)
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

func TestGetLeadTimeForChangesShouldRunWithInitialCommit(t *testing.T) {
	// https://github.com/go-git/go-git/releases/tag/v1.0.0
	commit1_0_0, err1 := repository.CommitObject(plumbing.NewHash("6f43e8933ba3c04072d5d104acc6118aac3e52ee"))
	if err1 != nil {
		t.Error(err1.Error())
	}

	duration := GetLeadTimeForChanges(repository, nil, commit1_0_0)
	expected, _ := time.ParseDuration("4152h5m44s")
	if duration.Minutes() == expected.Minutes() {
		return
	}
	t.Error(duration, "should be equal to ", expected)
}

func TestGetLeadTimeForChangesShouldReturnNilWhenNewerCommitIsNil(t *testing.T) {
	duration := GetLeadTimeForChanges(repository, nil, nil)
	if duration == nil {
		return
	}
	t.Error("should be equal to nil but", duration)
}
