package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"strings"
	"testing"

	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

func TestGetCommandReleaseShouldReturnReleasesWithRepositoryUrl(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git"})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	var releases []releases.Release
	json.Unmarshal(output.Bytes(), &releases)
	expectedReleasesNum := 60
	if len(releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(releases))
	}
}

func TestGetCommandReleaseShouldReturnReleasesWithRepositoryUrlSinceUntil(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	app := &cli.App{Writer: output}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"releases", "--repository", "https://github.com/go-git/go-git", "--since", "2020-01-01", "--until", "2020-12-31"})

	cCtx := cli.NewContext(app, set, nil)
	GetCommandReleases().Run(cCtx)

	var releases []releases.Release
	json.Unmarshal(output.Bytes(), &releases)
	expectedReleasesNum := 3
	if len(releases) != expectedReleasesNum {
		t.Errorf("releases should have %v releases but %v", expectedReleasesNum, len(releases))
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
