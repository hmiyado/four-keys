package main

import (
	"log"
	"os"

	"github.com/hmiyado/four-keys/internal/cli"
)

var version string

func main() {
	app := cli.DefaultApp(version)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
