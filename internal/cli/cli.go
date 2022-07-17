package cli

import (
	"encoding/json"
	"time"

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
	Option              *releases.Option         `json:"option"`
	DeploymentFrequency float64                  `json:"deploymentFrequency"`
	LeadTimeForChanges  LeadTimeForChangesOutput `json:"leadTimeForChanges"`
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
		LeadTimeForChanges:  getLeadTimeForChangesOutput(getMeanLeadTimeForChanges(output)),
	})
	if err != nil {
		context.Error(err)
		return err
	}
	context.Write(outputJson)
	return nil

}

func getMeanLeadTimeForChanges(output *ReleasesCliOutput) time.Duration {
	if len(output.Releases) == 0 {
		return time.Duration(0)
	}
	sum := time.Duration(0)
	for _, release := range output.Releases {
		sum = release.LeadTimeForChanges + sum
	}
	return time.Duration(int64(sum) / int64(len(output.Releases)))
}
