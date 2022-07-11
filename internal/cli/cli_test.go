package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestDefaultAppShouldReturnMetricsWithRepositoryUrlSinceUntil(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp()
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
	// {"option":{"since":"2020-01-01T00:00:00Z","until":"2020-12-31T23:59:59Z"},"deploymentFrequency":0.00821917808219178}
	if cliOutput.DeploymentFrequency < 0.00821 || cliOutput.DeploymentFrequency > 0.00822 {
		t.Errorf("deploymentFrequency should be in (0.00821, 0.00822) but %v", cliOutput.DeploymentFrequency)
	}
}

func TestDefaultAppShouldRunWithoutOption(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	defaltApp := DefaultApp()
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
