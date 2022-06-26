package main

import (
	"os"
	"log"
    "fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "boom",
        Usage: "make an explosive entrance",
        Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			repository, _ := git.PlainOpen("./")
			iter, _ := repository.CommitObjects()
			iter.ForEach(func (commit *object.Commit) error {
				fmt.Println(commit)
				return nil
			}); 

			return nil
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }

}
