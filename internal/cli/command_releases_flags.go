package cli

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hmiyado/four-keys/internal/core"
	"github.com/urfave/cli/v2"
)

func getCommandReleasesFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "repository",
			Usage:       "the remote repository url. repository will be cloned in memory",
			DefaultText: "local repository of current directory",
		},
		&cli.StringFlag{
			Name:        "accessToken",
			Usage:       "GitHub access token to clone private repository",
			DefaultText: "no access token",
		},
		&cli.TimestampFlag{
			Name:        "since",
			Usage:       "the start date to query releases (inclusive)",
			DefaultText: "1 month ago",
			Layout:      "2006-01-02",
		},
		&cli.TimestampFlag{
			Name:        "until",
			Usage:       "the end date to query releases (inclusive)",
			DefaultText: "now",
			Layout:      "2006-01-02",
		},
		&cli.StringFlag{
			Name:  "ignorePattern",
			Usage: "ignore releases that matches the pattern(regex)",
		},
		&cli.StringFlag{
			Name:        "fixCommitPattern",
			Usage:       "commit that message matches fixCommitPattern is regarded fix commit",
			DefaultText: "hotfix",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug message",
		},
	}
}

func onUsageError(cCtx *cli.Context, err error, isSubcommand bool) error {
	return err
}

var timerMap map[string]time.Time

func (c *CliContextWrapper) Since() time.Time {
	optionSince := c.context.Timestamp("since")
	if optionSince != nil {
		return *optionSince
	} else {
		now := time.Now()
		return time.Date(now.Year(), now.Month()-1, now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	}
}

func (c *CliContextWrapper) Until() time.Time {
	optionSince := c.context.Timestamp("until")
	if optionSince != nil {
		return (*optionSince).AddDate(0, 0, 1).Add(-time.Second)
	} else {
		return time.Now()
	}
}

func (c *CliContextWrapper) IgnorePattern() (*regexp.Regexp, error) {
	pattern := c.context.String("ignorePattern")
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

func (c *CliContextWrapper) FixCommitPattern() (*regexp.Regexp, error) {
	pattern := c.context.String("fixCommitPattern")
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

func (c *CliContextWrapper) Repository() (*git.Repository, error) {
	c.StartTimer("Open Repository")
	defer c.StopTimer("Open Repository")
	repositoryUrl := c.context.String("repository")
	accessToken := c.context.String("accessToken")
	var repository *git.Repository
	var error error
	var auth *http.BasicAuth
	if accessToken != "" {
		auth = &http.BasicAuth{
			Username: "four-keys",
			Password: accessToken,
		}
	}
	if repositoryUrl == "" {
		repository, error = git.PlainOpenWithOptions("./", &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: false})
		if error != nil {
			return nil, errors.New("cannot open repository at current directory")
		}
	} else {
		repository, error = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			Auth: auth,
			URL:  repositoryUrl,
		})
		if error != nil {
			return nil, fmt.Errorf("cannot clone repository: %v", repositoryUrl)
		}
	}

	return repository, nil
}

func (c *CliContextWrapper) Option() (*core.Option, error) {
	ignorePattern, err := c.IgnorePattern()
	if err != nil {
		wrappedError := fmt.Errorf("[invalid ignorePattern] %v", err)
		c.Error(wrappedError)
		return nil, wrappedError
	}

	fixCommitPattern, err := c.FixCommitPattern()
	if err != nil {
		wrappedError := fmt.Errorf("[invalid fixCommitPattern] %v", err)
		c.Error(wrappedError)
		return nil, wrappedError
	}

	return &core.Option{
		Since:             c.Since(),
		Until:             c.Until(),
		IgnorePattern:     ignorePattern,
		IsLocalRepository: c.context.String("repository") == "",
		FixCommitPattern:  fixCommitPattern,
		StartTimerFunc:    c.StartTimer,
		StopTimerFunc:     c.StopTimer,
		DebuglnFunc:       c.Debugln,
	}, nil
}
