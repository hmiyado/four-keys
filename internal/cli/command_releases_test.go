package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"strings"
	"testing"

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

func TestGetCommandReleaseShouldReturnReleasesWithRepositoryUrl(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git"})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	var cliOutput ReleasesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	expectedReleasesNum := 60
	if len(cliOutput.Releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(cliOutput.Releases))
	}
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
	_ = set.Parse([]string{
		"releases",
		"--repository", "https://github.com/go-git/go-git",
		"--since", "2020-01-01",
		"--until", "2020-12-31",
	})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	// Output should be
	// { "option":{"since":"2020-01-01T00:00:00Z","until":"2020-12-31T23:59:59Z","ignorePattern":null},
	//   "releases":[
	//     {"tag":"v5.2.0","date":"2020-10-09T11:49:30+02:00","leadTimeForChanges":{"value":130.77916666666667,"unit":"day"},"result":{"isSuccess":true}},
	//     {"tag":"v5.1.0","date":"2020-05-24T19:25:08+02:00","leadTimeForChanges":{"value":66.9150462962963,"unit":"day"},"result":{"isSuccess":true}},
	//     {"tag":"v5.0.0","date":"2020-03-15T21:18:32+01:00","leadTimeForChanges":{"value":224.73468749999998,"unit":"day"},"result":{"isSuccess":true}}]}⏎
	var cliOutput ReleasesCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	expectedReleasesNum := 3
	if len(cliOutput.Releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(cliOutput.Releases))
	}
	util.AssertIsNearBy(t, cliOutput.Releases[0].LeadTimeForChanges.Present(), 130.77916666666667, 0.01)
	util.AssertIsNearBy(t, cliOutput.Releases[1].LeadTimeForChanges.Present(), 66.9150462962963, 0.01)
	util.AssertIsNearBy(t, cliOutput.Releases[2].LeadTimeForChanges.Present(), 224.73468749999998, 0.01)
}

func TestGetCommandReleaseShouldBeFailWithInvalidSince(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git", "--since", "invalidtext"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx)

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
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git", "--until", "invalidtext"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx)

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
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git", "--ignorePattern", "*"})

	cCtx := cli.NewContext(app, set, nil)
	error := GetCommandReleases().Run(cCtx)

	if error == nil {
		t.Errorf("Invalid --until option does not return error. log: %v", output.String())
	}
	if !strings.Contains(error.Error(), "[invalid ignore pattern]") {
		t.Errorf("Invalid --until option does not return error of --until. error: %v", error.Error())
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
