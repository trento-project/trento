package runner

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/runner"
)

type RunnerCmdTestSuite struct {
	suite.Suite
	cmd *cobra.Command
}

func TestRunnerCmdTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerCmdTestSuite))
}

func (suite *RunnerCmdTestSuite) SetupTest() {
	os.Clearenv()

	cmd := NewRunnerCmd()
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// do nothing
	}

	cmd.Commands()[0].Run = func(cmd *cobra.Command, args []string) {
		// do nothing
	}

	cmd.SetArgs([]string{
		"start",
	})

	var b bytes.Buffer
	cmd.SetOut(&b)

	suite.cmd = cmd
}

func (suite *RunnerCmdTestSuite) TearDownTest() {
	suite.cmd.Execute()

	expectedConfig := &runner.Config{
		ApiHost:       "some-api-host",
		ApiPort:       1337,
		Interval:      1 * time.Minute,
		AnsibleFolder: "path/to/ansible",
	}
	config := LoadConfig()

	suite.EqualValues(expectedConfig, config)
}

func (suite *RunnerCmdTestSuite) TestConfigFromFlags() {
	suite.cmd.SetArgs([]string{
		"start",
		"--api-host=some-api-host",
		"--api-port=1337",
		"--interval=1",
		"--ansible-folder=path/to/ansible",
	})
}

func (suite *RunnerCmdTestSuite) TestConfigFromEnv() {
	os.Setenv("TRENTO_API_HOST", "some-api-host")
	os.Setenv("TRENTO_API_PORT", "1337")
	os.Setenv("TRENTO_INTERVAL", "1")
	os.Setenv("TRENTO_ANSIBLE_FOLDER", "path/to/ansible")
}

func (suite *RunnerCmdTestSuite) TestConfigFromFile() {
	os.Setenv("TRENTO_CONFIG", "../../test/fixtures/config/runner.yaml")
}
