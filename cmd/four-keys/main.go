package main

import (
	"os"
	"log"

    "github.com/hmiyado/four-keys/internal/cli"
)

func main() {
    app := cli.DefaultApp()

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }

}
