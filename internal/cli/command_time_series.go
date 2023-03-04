package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hmiyado/four-keys/internal/core"
	"github.com/urfave/cli/v2"
)

type TimeSeriesCliOutput struct {
	Option *core.Option          `json:"option"`
	Items  []TimeSeriesDataPoint `json:"items"`
}

type TimeSeriesDataPoint struct {
	Date                time.Time `json:"time"`
	DeploymentFrequency float64   `json:"deploymentFrequency"`
	LeadTimeForChanges  float64   `json:"leadTimeForChanges"`
	TimeToRestore       float64   `json:"timeToRestore"`
	ChangeFailureRate   float64   `json:"changeFailureRate"`
}

func GetCommandTimeSeries() *cli.Command {
	return &cli.Command{
		Name:  "timeSeries",
		Usage: "list time series of four keys",
		Flags: append(getCommandReleasesFlags(), getCommandTimeSeriesFlags()...),
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
			timeSeriesOption, err := context.TimeSeriesOption()
			if err != nil {
				return err
			}

			if !validateTimeSeriesInterval(timeSeriesOption.Interval, option.Since, option.Until) {
				err := fmt.Errorf("Interval is too short")
				return err
			}
			output := &TimeSeriesCliOutput{
				Option: option,
				Items:  mapReleasesToTimeSeriesCliOutput(releases, timeSeriesOption.Interval, option.Since, option.Until),
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

func validateTimeSeriesInterval(interval TimeSeriesInterval, since time.Time, until time.Time) bool {
	duration := until.Sub(since)
	switch interval {
	case Day:
		if duration < 24*time.Hour {
			return false
		}
	case Week:
		if duration < 7*24*time.Hour {
			return false
		}
	case Month:
		if duration < 28*24*time.Hour {
			return false
		}
	}
	return true
}

func mapReleasesToTimeSeriesCliOutput(releases []*core.Release, interval TimeSeriesInterval, since time.Time, until time.Time) []TimeSeriesDataPoint {
	var items []TimeSeriesDataPoint
	dateOfStart := until
	switch interval {
	case Week:
		dateOfStart = time.Date(until.Year(), until.Month(), until.Day()-int(until.Weekday()), 0, 0, 0, 0, until.Location())
	case Month:
		dateOfStart = time.Date(until.Year(), until.Month(), 1, 0, 0, 0, 0, until.Location())
	}
	dateOfEnd := until
	for ; dateOfStart.After(since) || dateOfStart.Equal(since); dateOfStart = getBeforeDate(dateOfStart, interval) {
		items = append(items, mapReleasesToTimeSeriesDataPoint(releases, dateOfStart, dateOfEnd, interval))
		dateOfEnd = dateOfStart
	}
	return items
}

func mapReleasesToTimeSeriesDataPoint(releases []*core.Release, dateOfStart time.Time, dateOfEnd time.Time, interval TimeSeriesInterval) TimeSeriesDataPoint {
	var releasesInInterval []*core.Release
	for _, release := range releases {
		if release.Date.After(dateOfStart) && release.Date.Before(dateOfEnd) {
			releasesInInterval = append(releasesInInterval, release)
		}
	}
	return TimeSeriesDataPoint{
		Date:                dateOfStart,
		DeploymentFrequency: core.GetDeploymentFrequencyByTimeunit(releasesInInterval, dateOfStart, dateOfEnd, string(interval)),
		LeadTimeForChanges:  core.GetMeanLeadTimeForChanges(releasesInInterval).Hours(),
		TimeToRestore:       core.GetTimeToRestore(releasesInInterval).Hours(),
		ChangeFailureRate:   core.GetChangeFailureRate(releasesInInterval),
	}
}

func getBeforeDate(date time.Time, interval TimeSeriesInterval) time.Time {
	return getNextDate(date, interval, true)
}

func getNextDate(date time.Time, interval TimeSeriesInterval, goBefore bool) time.Time {
	sign := 1
	if goBefore {
		sign = -1
	}

	switch interval {
	case Day:
		return date.AddDate(0, 0, 1*sign)
	case Week:
		return date.AddDate(0, 0, 7*sign)
	case Month:
		return date.AddDate(0, 1*sign, 0)
	}
	return date
}
