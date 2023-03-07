package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type TimeSeriesInterval string

const (
	Day   TimeSeriesInterval = "day"
	Week  TimeSeriesInterval = "week"
	Month TimeSeriesInterval = "month"
)

type TimeSeriesOption struct {
	Interval TimeSeriesInterval
}

func (c *CliContextWrapper) Interval() (*TimeSeriesInterval, error) {
	intervalString := c.context.String("interval")
	if intervalString == "" {
		interval := Month
		return &interval, nil
	}
	validIntervals := []TimeSeriesInterval{Day, Week, Month}
	for _, interval := range validIntervals {
		if intervalString == string(interval) {
			return &interval, nil
		}
	}
	return nil, fmt.Errorf("unavailable interval \"%s\". Interval should be one of %s", intervalString, validIntervals)
}

func (c *CliContextWrapper) TimeSeriesOption() (*TimeSeriesOption, error) {
	interval, error := c.Interval()
	if error != nil {
		return nil, error
	}
	return &TimeSeriesOption{
		Interval: *interval,
	}, nil
}

func getCommandTimeSeriesFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "interval",
			Usage:       "Interval of time series: day, week, month",
			DefaultText: "month",
		},
	}
}
