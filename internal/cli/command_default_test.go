package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hmiyado/four-keys/internal/util"
	"github.com/urfave/cli/v2"
)

func TestDefaultAppShouldReturnMetricsWithRepositoryUrlSinceUntil(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp("")
	testApp := &cli.App{
		Flags:  defaltApp.Flags,
		Action: defaltApp.Action,
		Writer: output,
	}

	err := testApp.Run([]string{"four-keys", "--repository", "https://github.com/go-git/go-git", "--since", "2020-01-01", "--until", "2020-12-31"})
	if err != nil {
		t.Errorf(err.Error())
	}
	var cliOutput DefaultCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	// intended output
	// {
	//   "option":{"since":"2020-01-01T00:00:00Z","until":"2020-12-31T23:59:59Z"},
	//   "deploymentFrequency":0.00821917808219178,
	//   "leadTimeForChanges":{"value":140.8096334876543,"unit":"day"},
	//   "changeFailureRate":0
	// }
	util.AssertIsNearBy(t, cliOutput.DeploymentFrequency, 0.00821917808219178, 0.01)
	util.AssertIsNearBy(t, cliOutput.LeadTimeForChanges.Present(), 140.8096334876543, 0.01)
	util.AssertIsNearBy(t, cliOutput.TimeToRestore.Present(), 0, 0.01)
	util.AssertIsNearBy(t, cliOutput.ChangeFailureRate, 0, 0.01)
}

func TestDefaultAppShouldReturnTimeToRestoreAndChangeFailureRate(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp("")
	testApp := &cli.App{
		Flags:  defaltApp.Flags,
		Action: defaltApp.Action,
		Writer: output,
	}

	err := testApp.Run([]string{"four-keys", "--repository", "https://github.com/go-git/go-git", "--since", "2015-12-01", "--until", "2016-05-31"})
	if err != nil {
		t.Errorf(err.Error())
	}
	var cliOutput DefaultCliOutput
	json.Unmarshal(output.Bytes(), &cliOutput)
	util.AssertIsNearBy(t, cliOutput.ChangeFailureRate, 0.08333333333333333, 0.01)
	util.AssertIsNearBy(t, cliOutput.TimeToRestore.Present(), 2.7969791666666666, 0.01)
}

func TestDefaultAppShouldRunWithoutOption(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp("")
	testApp := &cli.App{
		Flags:  defaltApp.Flags,
		Action: defaltApp.Action,
		Writer: output,
	}

	err := testApp.Run([]string{"four-keys"})
	if err != nil {
		t.Errorf(err.Error())
	}

}
