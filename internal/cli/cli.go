package cli

import (
	"encoding/json"

	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

func DefaultApp() *cli.App {
	return &cli.App{
		Name:   "four-keys",
		Usage:  "analyze four keys metrics",
		Flags:  getCommandReleasesFlags(),
		Action: defaultAction,
		Commands: []*cli.Command{
			GetCommandReleases(),
		},
	}
}

type DefaultCliOutput struct {
	Option              *releases.Option `json:"option"`
	DeploymentFrequency float64          `json:"deploymentFrequency"`
}

func defaultAction(ctx *cli.Context) error {
	context := &CliContextWrapper{context: ctx}
	output, err := QueryReleases(context)
	if err != nil {
		context.Error(err)
		return err
	}

	duration := output.Option.Until.Sub(output.Option.Since)
	daysCount := int(duration.Hours() / 24)
	releasesCount := len(output.Releases)
	deploymentFrequency := float64(releasesCount) / float64(daysCount)

	outputJson, err := json.Marshal(&DefaultCliOutput{
		Option:              output.Option,
		DeploymentFrequency: deploymentFrequency,
	})
	if err != nil {
		context.Error(err)
		return err
	}
	context.Write(outputJson)
	return nil

}
