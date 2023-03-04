package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
)

func TestGetCommandTimeSeriesShouldReturnTimeSeriesWithoutOptions(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"timeSeries"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandTimeSeries().Run(cCtx)

	if error != nil {
		t.Errorf("timeSeries without options failed. error:%v", error.Error())
	}

}

func TestGetCommandTimeSeriesShouldReturnTimeSeriesInDayly(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	args := []string{
		"timeSeries",
		"--repository", "https://github.com/hmiyado/four-keys",
		"--since", "2022-10-01",
		"--until", "2022-10-31",
		"--interval", "day"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	GetCommandTimeSeries().Run(cCtx, args...)

	var cliOutput TimeSeriesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)

	first := cliOutput.Items[0]
	if first.Date.Sub(cliOutput.Items[1].Date) != time.Hour*24 {
		t.Errorf("timeSeries should be dayly but %v", first.Date.Sub(cliOutput.Items[1].Date))
	}
	if first.Date.Year() != 2022 || first.Date.Month() != 10 || first.Date.Day() != 31 {
		t.Errorf("timeSeries should start from 2022-10-31 but %v", first.Date)
	}
	last := cliOutput.Items[len(cliOutput.Items)-1]
	if last.Date.Year() != 2022 || last.Date.Month() != 10 || last.Date.Day() != 1 {
		t.Errorf("timeSeries should end at 2022-10-01 but %v", last.Date)
	}
}

func TestGetCommandTimeSeriesShouldReturnTimeSeriesInWeekly(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	args := []string{
		"timeSeries",
		"--repository", "https://github.com/hmiyado/four-keys",
		"--since", "2022-10-01",
		"--until", "2022-10-31",
		"--interval", "week"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	GetCommandTimeSeries().Run(cCtx, args...)

	var cliOutput TimeSeriesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)

	first := cliOutput.Items[0]
	if first.Date.Sub(cliOutput.Items[1].Date) != time.Hour*24*7 {
		t.Errorf("timeSeries should be weekly but %v", first.Date.Sub(cliOutput.Items[1].Date))
	}
	// 2022-10-30 is Sunday
	if first.Date.Year() != 2022 || first.Date.Month() != 10 || first.Date.Day() != 30 {
		t.Errorf("timeSeries should start from 2022-10-30 but %v", first.Date)
	}
}

func TestGetCommandTimeSeriesShouldReturnTimeSeriesInMonthly(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	args := []string{
		"timeSeries",
		"--repository", "https://github.com/hmiyado/four-keys",
		"--since", "2022-09-01",
		"--until", "2022-10-31",
		"--interval", "month"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	GetCommandTimeSeries().Run(cCtx, args...)

	var cliOutput TimeSeriesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)

	first := cliOutput.Items[0]
	if first.Date.Month() == cliOutput.Items[1].Date.Month()-1 {
		t.Errorf("timeSeries should be monthly but %v and %v", first.Date, cliOutput.Items[1].Date)
	}
	if first.Date.Year() != 2022 || first.Date.Month() != 10 || first.Date.Day() != 1 {
		t.Errorf("timeSeries should start from 2022-10-01 but %v", first.Date)
	}
}
