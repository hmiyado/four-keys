package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

func DefaultApp() *cli.App {
	return &cli.App{
		Name:  "four-keys",
		Usage: "analyze four keys metrics",
		Action: func(*cli.Context) error {
			repository, _ := git.PlainOpen("./")
			iter, _ := repository.CommitObjects()
			iter.ForEach(func(commit *object.Commit) error {
				fmt.Println(commit)
				return nil
			})

			return nil
		},
		Commands: []*cli.Command{
			getCommandReleases(),
		},
	}
}

func getCommandReleases() *cli.Command {
	return &cli.Command{
		Name:  "releases",
		Usage: "list releases",
		Flags: []cli.Flag{
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
		},
		Action: func(ctx *cli.Context) error {
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

			repositoryUrl := ctx.String("repository")
			var repository *git.Repository
			var error error
			if repositoryUrl == "" {
				repository, error = git.PlainOpen("./")
				if error != nil {
					fmt.Errorf("cannot open repository at current directory")
					fmt.Errorf(error.Error())
					return error
				}
			} else {
				repository, error = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
					URL: repositoryUrl,
				})
				if error != nil {
					fmt.Errorf("cannot clone repository: %v", repositoryUrl)
					fmt.Errorf(error.Error())
					return error
				}
			}

			releases := releases.QueryReleases(repository, &releases.Option{
				StartDate: since,
				EndDate:   until,
			})
			releasesJson, error := json.Marshal(releases)
			if error != nil {
				fmt.Errorf(error.Error())
				return error
			}
			fmt.Printf("%s", releasesJson)
			return nil
		},
	}
}
