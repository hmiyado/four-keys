package cli

import (
	"encoding/json"

	"github.com/hmiyado/four-keys/internal/core"
	"github.com/urfave/cli/v2"
)

type DefaultCliOutput struct {
	Option              *core.Option         `json:"option"`
	DeploymentFrequency float64              `json:"deploymentFrequency"`
	LeadTimeForChanges  DurationWithTimeUnit `json:"leadTimeForChanges"`
	TimeToRestore       DurationWithTimeUnit `json:"timeToRestore"`
	ChangeFailureRate   float64              `json:"changeFailureRate"`
}

func defaultAction(ctx *cli.Context) error {
	context := &CliContextWrapper{context: ctx}
	context.Debugln("In debug mode")
	releases, err := QueryReleases(context)
	if err != nil {
		context.Error(err)
		return err
	}
	option, err := context.Option()
	if err != nil {
		context.Error(err)
		return err
	}

	context.StartTimer("Calculate metrics")
	outputJson, err := json.Marshal(&DefaultCliOutput{
		Option:              option,
		DeploymentFrequency: getDeploymentFrequency(releases, *option),
		LeadTimeForChanges:  getDurationWithTimeUnit(getMeanLeadTimeForChanges(releases)),
		TimeToRestore:       getDurationWithTimeUnit(getTimeToRestore(releases)),
		ChangeFailureRate:   getChangeFailureRate(releases),
	})
	context.StopTimer("Calculate metrics")
	if err != nil {
		context.Error(err)
		return err
	}
	context.Write(outputJson)
	return nil

}
