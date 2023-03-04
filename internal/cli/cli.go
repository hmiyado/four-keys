package cli

import (
	"github.com/urfave/cli/v2"
)

func DefaultApp(version string) *cli.App {
	return &cli.App{
		Name:    "four-keys",
		Usage:   "analyze four keys metrics",
		Version: version,
		Flags:   getCommandReleasesFlags(),
		Action:  defaultAction,
		Commands: []*cli.Command{
			GetCommandReleases(),
		},
		OnUsageError: onUsageError,
	}
}
