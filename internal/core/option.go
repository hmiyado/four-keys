package core

import (
	"regexp"
	"strings"
	"time"
)

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
