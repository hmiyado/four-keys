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
}

type ReleasesCliOutput struct {
	Option   *releases.Option    `json:"option"`
	Releases []*releases.Release `json:"releases"`
}

func GetCommandReleases() *cli.Command {
	return &cli.Command{
		Name:  "releases",
		Usage: "list releases",
		Flags: CommandReleasesFlags,
		Action: func(ctx *cli.Context) error {
			output, err := QueryReleases(ctx)

			if err != nil {
				ctx.App.ErrWriter.Write([]byte(err.Error()))
				return err
			}
			releasesJson, error := json.Marshal(output)
			if error != nil {
				ctx.App.ErrWriter.Write([]byte(error.Error()))
				return error
			}
			ctx.App.Writer.Write(releasesJson)
			return nil
		},
	}
}

func parseOptionSinceUntil(ctx *cli.Context) (time.Time, time.Time) {
	optionSince := ctx.Timestamp("since")
	optionUntil := ctx.Timestamp("until")
	var since, until time.Time = time.Unix(0, 0), time.Now()
	if optionSince != nil {
		since = *optionSince
	}
	if optionUntil != nil {
		until = *optionUntil
		until = until.AddDate(0, 0, 1).Add(-time.Second)
	}

	return since, until
}

func parseOptionRepository(ctx *cli.Context) (*git.Repository, error) {
	repositoryUrl := ctx.String("repository")
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

func QueryReleases(ctx *cli.Context) (*ReleasesCliOutput, error) {
	since, until := parseOptionSinceUntil(ctx)

	repository, error := parseOptionRepository(ctx)
	if error != nil {
		ctx.App.ErrWriter.Write([]byte(error.Error()))
		return nil, error
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
