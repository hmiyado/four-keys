package core

import (
	"fmt"
	"time"
)

type Release struct {
	Tag                string        `json:"tag"`
	Date               time.Time     `json:"date"`
	LeadTimeForChanges time.Duration `json:"leadTimeForChanges"`
	Result             ReleaseResult `json:"result"`
	isRestored         bool          `json:"-"`
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
