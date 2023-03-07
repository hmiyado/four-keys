package cli

import (
	"encoding/json"
	"time"

	"github.com/hmiyado/four-keys/internal/core"
	"github.com/urfave/cli/v2"
)

type ReleasesCliOutput struct {
	Option   *core.Option        `json:"option"`
	Releases []*ReleaseCliOutput `json:"releases"`
}

type ReleaseCliOutput struct {
	Tag                string                 `json:"tag"`
	Date               time.Time              `json:"date"`
	LeadTimeForChanges DurationWithTimeUnit   `json:"leadTimeForChanges"`
	Result             ReleaseResultCliOutput `json:"result"`
}

type ReleaseResultCliOutput struct {
	IsSuccess     bool                  `json:"isSuccess"`
	TimeToRestore *DurationWithTimeUnit `json:"timeToRestore"`
}

func GetCommandReleases() *cli.Command {
	return &cli.Command{
		Name:  "releases",
		Usage: "list releases",
		Flags: getCommandReleasesFlags(),
		Action: func(ctx *cli.Context) error {
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

			output := &ReleasesCliOutput{
				Option:   option,
				Releases: mapReleasesToCliOutput(releases),
			}
			releasesJson, err := json.Marshal(output)
			if err != nil {
				context.Error(err)
				return err
			}
			context.Write(releasesJson)
			return nil
		},
		OnUsageError: onUsageError,
	}
}

func QueryReleases(context *CliContextWrapper) ([]*core.Release, error) {
	repository, err := context.Repository()
	if err != nil {
		context.Error(err)
		return nil, err
	}

	option, err := context.Option()
	if err != nil {
		return nil, err
	}
	return core.QueryReleases(repository, option), nil
}

func mapReleasesToCliOutput(releases []*core.Release) []*ReleaseCliOutput {
	output := make([]*ReleaseCliOutput, 0)
	for _, release := range releases {
		output = append(output, &ReleaseCliOutput{
			Tag:                release.Tag,
			Date:               release.Date,
			LeadTimeForChanges: getDurationWithTimeUnit(release.LeadTimeForChanges),
			Result:             mapReleaseResultToCliOutput(release.Result),
		})
	}
	return output
}

func mapReleaseResultToCliOutput(result core.ReleaseResult) ReleaseResultCliOutput {
	if result.TimeToRestore == nil {
		return ReleaseResultCliOutput{
			IsSuccess:     result.IsSuccess,
			TimeToRestore: nil,
		}
	}
	t := getDurationWithTimeUnit(*result.TimeToRestore)
	return ReleaseResultCliOutput{
		IsSuccess:     result.IsSuccess,
		TimeToRestore: &t,
	}
}
