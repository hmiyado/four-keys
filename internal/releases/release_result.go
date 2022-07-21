package releases

import (
	"fmt"
	"time"
)

type ReleaseResult struct {
	IsSuccess     bool
	TimeToRestore *time.Duration
}

func (r ReleaseResult) Equal(another ReleaseResult) bool {
	if r.TimeToRestore == nil && another.TimeToRestore == nil {
		return r.IsSuccess == another.IsSuccess
	}
	if r.TimeToRestore != nil && another.TimeToRestore != nil {
		return r.IsSuccess == another.IsSuccess && *(r.TimeToRestore) == *(another.TimeToRestore)
	}
	return false
}

func (r ReleaseResult) String() string {
	return fmt.Sprintf("IsSuccess=%v, TimeToRestore=%v", r.IsSuccess, r.TimeToRestore)
}
