package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hmiyado/four-keys/internal/util"
	"github.com/urfave/cli/v2"
)

func TestGetCommandReleaseShouldReturnReleasesWithoutOptions(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx)

	if error != nil {
		t.Errorf("releases without options failed. error:%v", error.Error())
	}
}

func TestGetCommandReleaseShouldHaveDefaultTimeRangeOption(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git"})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	var cliOutput ReleasesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	days28, _ := time.ParseDuration(fmt.Sprintf("%vh", 24*28))
	days31, _ := time.ParseDuration(fmt.Sprintf("%vh", 24*31))
	duration := cliOutput.Option.Until.Sub(cliOutput.Option.Since)
	if (duration-days28) > -time.Second && (duration-days31) < time.Second {
		return
	}
	t.Logf("option: %v", cliOutput.Option)
	t.Errorf("time range should have abount 1 month(28-31days) but %v", duration)
}

func TestGetCommandReleaseShouldReturnReleasesWithRepositoryUrlSinceUntilIgnorePattern(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{
		"releases",
		"--repository", "https://github.com/go-git/go-git",
		"--since", "2020-01-01",
		"--until", "2020-12-31",
		"--ignorePattern", "v5\\.2\\.0",
	})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	var cliOutput ReleasesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	expectedReleasesNum := 2
	if len(cliOutput.Releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(cliOutput.Releases))
	}
}

func TestGetCommandReleaseShouldHaveLeadTimeForChangesForEachReleases(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	args := []string{
		"releases",
		"--repository", "https://github.com/go-git/go-git",
		"--since", "2020-01-01",
		"--until", "2020-12-31",
	}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx, args...)

	// Output:
	// {
	//   "option": { "since": "2020-01-01T00:00:00Z", "until": "2020-12-31T23:59:59Z" },
	//   "releases": [
	//     {
	//       "tag": "v5.2.0",
	//       "date": "2020-10-09T11:49:30+02:00",
	//       "leadTimeForChanges": { "value": 130.77916666666667, "unit": "day" },
	//       "result": { "isSuccess": true, "timeToRestore": null }
	//     },
	//     {
	//       "tag": "v5.1.0",
	//       "date": "2020-05-24T19:25:08+02:00",
	//       "leadTimeForChanges": { "value": 69.86515046296296, "unit": "day" },
	//       "result": { "isSuccess": true, "timeToRestore": null }
	//     },
	//     {
	//       "tag": "v5.0.0",
	//       "date": "2020-03-15T21:18:32+01:00",
	//       "leadTimeForChanges": {
	//         "value": 224.73468749999998,
	//         "unit": "day"
	//       },
	//       "result": { "isSuccess": true, "timeToRestore": null }
	//     }
	//   ]
	// }
	var cliOutput ReleasesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	expectedReleasesNum := 3
	if len(cliOutput.Releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(cliOutput.Releases))
	}
	util.AssertIsNearBy(t, cliOutput.Releases[0].LeadTimeForChanges.Present(), 130.77916666666667, 0.01)
	util.AssertIsNearBy(t, cliOutput.Releases[1].LeadTimeForChanges.Present(), 69.86515046296296, 0.01)
	util.AssertIsNearBy(t, cliOutput.Releases[2].LeadTimeForChanges.Present(), 224.73468749999998, 0.01)
}

func TestGetCommandReleaseShouldBeFailWithInvalidSince(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	errOutput := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output, ErrWriter: errOutput}
	set := flag.NewFlagSet("test", 0)
	args := []string{"releases", "--repository", "https://github.com/go-git/go-git", "--since", "invalidtext"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx, args...)

	if error == nil {
		t.Errorf("Invalid --since option does not return error. log: %v", output.String())
	}
	if !strings.Contains(error.Error(), "since") {
		t.Errorf("Invalid --since option does not return error of --since. erro: %v", error.Error())
	}
}

func TestGetCommandReleaseShouldBeFailWithInvalidUntil(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	args := []string{"releases", "--repository", "https://github.com/go-git/go-git", "--until", "invalidtext"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx, args...)

	if error == nil {
		t.Errorf("Invalid --until option does not return error. log: %v", output.String())
	}
	if !strings.Contains(error.Error(), "until") {
		t.Errorf("Invalid --until option does not return error of --until. erro: %v", error.Error())
	}
}

func TestGetCommandReleaseShouldBeFailWithInvalidIgnorePattern(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	errOutput := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output, ErrWriter: errOutput}
	set := flag.NewFlagSet("test", 0)
	args := []string{"releases", "--repository", "https://github.com/go-git/go-git", "--ignorePattern", "*"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx, args...)

	if error == nil {
		t.Errorf("Invalid --ignorePattern option does not return error. log: %v", output.String())
	}
	if !strings.Contains(error.Error(), "[invalid ignorePattern]") {
		t.Errorf("Invalid --ignorePattern option does not return error of --until. error: %v", error.Error())
	}
}

func TestGetCommandReleaseShouldBeFailWithInvalidFixCommitPattern(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	errOutput := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output, ErrWriter: errOutput}
	set := flag.NewFlagSet("test", 0)
	args := []string{"releases", "--repository", "https://github.com/go-git/go-git", "--fixCommitPattern", "*"}
	_ = set.Parse(args)

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx, args...)

	if error == nil {
		t.Errorf("Invalid --fixCommitPattern option does not return error. log: %v", output.String())
	}
	if !strings.Contains(error.Error(), "[invalid fixCommitPattern]") {
		t.Errorf("Invalid --fixCommitPattern option does not return error of --until. error: %v", error.Error())
	}
}
func TestGetCommandReleaseShouldBeDebuggable(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--debug"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx)

	if error != nil {
		t.Errorf("--debug is not available log: %v", output.String())
	}
}
