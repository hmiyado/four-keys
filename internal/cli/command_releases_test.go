package cli

import (
	"bytes"
	"encoding/json"
	"flag"
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
