package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

var CommandReleasesFlags []cli.Flag = []cli.Flag{
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
	&cli.BoolFlag{
		Name:  "debug",
		Usage: "show debug message",
	},
}

type ReleasesCliOutput struct {
	Option   *releases.Option    `json:"option"`
	Releases []*releases.Release `json:"releases"`
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
		Flags: CommandReleasesFlags,
		Action: func(ctx *cli.Context) error {
			context := &CliContextWrapper{context: ctx}
			context.Debugln("In debug mode")
			output, err := QueryReleases(context)

			if err != nil {
				context.Error(err)
				return err
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

func QueryReleases(context *CliContextWrapper) (*ReleasesCliOutput, error) {
	since := context.Since()
	until := context.Until()
	repository, err := context.Repository()

	if err != nil {
		context.Error(err)
		return nil, err
	}

	option := &releases.Option{
		StartDate: since,
		EndDate:   until,
	}
	releases := releases.QueryReleases(repository, option)

	return &ReleasesCliOutput{
		Option:   option,
		Releases: releases,
	}, nil

}
