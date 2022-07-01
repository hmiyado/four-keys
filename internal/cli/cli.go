package cli

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
			GetCommandReleases(),
		},
	}
}
