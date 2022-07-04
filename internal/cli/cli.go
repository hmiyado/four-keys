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
		Flags:  CommandReleasesFlags,
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
	output, err := QueryReleases(ctx)
	if err != nil {
		ctx.App.ErrWriter.Write([]byte(err.Error()))
		return err
	}

	duration := output.Option.EndDate.Sub(output.Option.StartDate)
	daysCount := int(duration.Hours() / 24)
	releasesCount := len(output.Releases)
	deploymentFrequency := float64(releasesCount) / float64(daysCount)

	outputJson, error := json.Marshal(&DefaultCliOutput{
		Option:              output.Option,
		DeploymentFrequency: deploymentFrequency,
	})
	if error != nil {
		ctx.App.ErrWriter.Write([]byte(error.Error()))
		return error
	}
	ctx.App.Writer.Write(outputJson)
	return nil

}
