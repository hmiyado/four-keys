package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"strings"
	"testing"

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
