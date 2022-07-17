package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

func getCommandReleasesFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "repository",
			Usage: "the remote repository url. repository will be cloned in memory. default is local repository(current directory)",
		},
		&cli.TimestampFlag{
			Name:   "since",
			Usage:  "the start date to query releases (inclusive)",
			Layout: "2006-01-02",
		},
		&cli.TimestampFlag{
			Name:   "until",
			Usage:  "the end date to query releases (inclusive)",
			Layout: "2006-01-02",
		},
		&cli.StringFlag{
			Name:  "ignorePattern",
			Usage: "ignore releases that matches the pattern(regex)",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug message",
		},
	}
}

type ReleasesCliOutput struct {
	Option   *releases.Option    `json:"option"`
	Releases []*ReleaseCliOutput `json:"releases"`
}

type ReleaseCliOutput struct {
	Tag                string                   `json:"tag"`
	Date               time.Time                `json:"date"`
	LeadTimeForChanges LeadTimeForChangesOutput `json:"leadTimeForChanges"`
}

type CliContextWrapper struct {
	context *cli.Context
}

func (c *CliContextWrapper) Since() time.Time {
	optionSince := c.context.Timestamp("since")
	if optionSince != nil {
		return *optionSince
	} else {
		return time.Unix(0, 0)
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

func (c *CliContextWrapper) Repository() (*git.Repository, error) {
	repositoryUrl := c.context.String("repository")
	var repository *git.Repository
	var error error
	if repositoryUrl == "" {
		repository, error = git.PlainOpenWithOptions("./", &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: false})
		if error != nil {
			return nil, errors.New("cannot open repository at current directory")
		}
	} else {
		repository, error = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: repositoryUrl,
		})
		if error != nil {
			return nil, fmt.Errorf("cannot clone repository: %v", repositoryUrl)
		}
	}

	return repository, nil
}

func (c *CliContextWrapper) Option() (*releases.Option, error) {
	ignorePattern, err := c.IgnorePattern()
	if err != nil {
		wrappedError := fmt.Errorf("[invalid ignore pattern] %v", err)
		c.Error(wrappedError)
		return nil, wrappedError
	}

	return &releases.Option{
		Since:         c.Since(),
		Until:         c.Until(),
		IgnorePattern: ignorePattern,
	}, nil
}

func (c *CliContextWrapper) Debugf(format string, a ...any) {
	debug := c.context.Bool("debug")
	if debug {
		fmt.Print("[Debug] ")
		fmt.Printf(format, a...)
	}
}

func (c *CliContextWrapper) Debugln(a ...any) {
	debug := c.context.Bool("debug")
	if debug {
		fmt.Print("[Debug] ")
		fmt.Println(a...)
	}
}

func (c *CliContextWrapper) Error(err error) {
	c.context.App.ErrWriter.Write([]byte(err.Error()))
}

func (c *CliContextWrapper) Write(p []byte) {
	c.context.App.Writer.Write(p)
}

func GetCommandReleases() *cli.Command {
	return &cli.Command{
		Name:  "releases",
		Usage: "list releases",
		Flags: getCommandReleasesFlags(),
		Action: func(ctx *cli.Context) error {
			context := &CliContextWrapper{context: ctx}
			context.Debugln("In debug mode")
			releases, err := QueryReleases(context)
			if err != nil {
				context.Error(err)
				return err
			}
			option, err := context.Option()
			if err != nil {
				context.Error(err)
				return err
			}

			output := &ReleasesCliOutput{
				Option:   option,
				Releases: mapReleasesToCliOutput(releases),
			}
			releasesJson, err := json.Marshal(output)
			if err != nil {
				context.Error(err)
				return err
			}
			context.Write(releasesJson)
			return nil
		},
	}
}

func QueryReleases(context *CliContextWrapper) ([]*releases.Release, error) {
	repository, err := context.Repository()
	if err != nil {
		context.Error(err)
		return nil, err
	}

	option, err := context.Option()
	if err != nil {
		return nil, err
	}
	return releases.QueryReleases(repository, option), nil
}

func mapReleasesToCliOutput(releases []*releases.Release) []*ReleaseCliOutput {
	output := make([]*ReleaseCliOutput, 0)
	for _, release := range releases {
		output = append(output, &ReleaseCliOutput{
			Tag:                release.Tag,
			Date:               release.Date,
			LeadTimeForChanges: getLeadTimeForChangesOutput(release.LeadTimeForChanges),
		})
	}
	return output
}
