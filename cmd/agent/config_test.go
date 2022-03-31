package agent

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/agent/discovery"
	"github.com/trento-project/trento/agent/discovery/collector"
)

type AgentCmdTestSuite struct {
	suite.Suite
	cmd *cobra.Command
}

func TestAgentCmdTestSuite(t *testing.T) {
	suite.Run(t, new(AgentCmdTestSuite))
}

func (suite *AgentCmdTestSuite) SetupTest() {
	os.Clearenv()

	cmd := NewAgentCmd()

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

func (suite *AgentCmdTestSuite) TearDownTest() {
	suite.cmd.Execute()

	expectedConfig := &agent.Config{
		InstanceName: "some-hostname",
		DiscoveriesConfig: &discovery.DiscoveriesConfig{
			SSHAddress: "some-ssh-address",
			DiscoveriesPeriodsConfig: &discovery.DiscoveriesPeriodConfig{
				Cluster:      10 * time.Second,
				SAPSystem:    10 * time.Second,
				Cloud:        10 * time.Second,
				Host:         10 * time.Second,
				Subscription: 900 * time.Second,
			},
			CollectorConfig: &collector.Config{
				CollectorHost: "localhost",
				CollectorPort: 1337,
				EnablemTLS:    true,
				Cert:          "some-cert",
				Key:           "some-key",
				CA:            "some-ca",
			},
		},
	}

	config, err := LoadConfig()
	config.InstanceName = "some-hostname"
	suite.NoError(err)

	suite.EqualValues(expectedConfig, config)
}

func (suite *AgentCmdTestSuite) TestConfigFromFlags() {
	suite.cmd.SetArgs([]string{
		"start",
		"--ssh-address=some-ssh-address",
		"--discovery-period=10",
		"--collector-host=localhost",
		"--collector-port=1337",
		"--enable-mtls",
		"--cert=some-cert",
		"--key=some-key",
		"--ca=some-ca",
	})
}

func (suite *AgentCmdTestSuite) TestConfigFromEnv() {
	os.Setenv("TRENTO_SSH_ADDRESS", "some-ssh-address")
	os.Setenv("TRENTO_DISCOVERY_PERIOD", "10")
	os.Setenv("TRENTO_COLLECTOR_HOST", "localhost")
	os.Setenv("TRENTO_COLLECTOR_PORT", "1337")
	os.Setenv("TRENTO_ENABLE_MTLS", "true")
	os.Setenv("TRENTO_CERT", "some-cert")
	os.Setenv("TRENTO_KEY", "some-key")
	os.Setenv("TRENTO_CA", "some-ca")
}

func (suite *AgentCmdTestSuite) TestConfigFromFile() {
	os.Setenv("TRENTO_CONFIG", "../../test/fixtures/config/agent.yaml")
}
