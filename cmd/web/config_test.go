package web

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/internal/db"
	"github.com/trento-project/trento/web"
)

type WebCmdTestSuite struct {
	suite.Suite
	cmd *cobra.Command
}

func TestWebCmdTestSuite(t *testing.T) {
	suite.Run(t, new(WebCmdTestSuite))
}

func (suite *WebCmdTestSuite) SetupTest() {
	os.Clearenv()

	cmd := NewWebCmd()

	cmd.Commands()[1].Run = func(cmd *cobra.Command, args []string) {
		// do nothing
	}

	cmd.SetArgs([]string{
		"serve",
	})

	var b bytes.Buffer
	cmd.SetOut(&b)

	suite.cmd = cmd
}

func (suite *WebCmdTestSuite) TearDownTest() {
	suite.cmd.Execute()

	expectedConfig := &web.Config{
		Host:          "some-host",
		Port:          1337,
		CollectorPort: 1338,
		EnablemTLS:    true,
		Cert:          "some-cert",
		Key:           "some-key",
		CA:            "some-ca",
		DBConfig: &db.Config{
			Host:     "some-db-host",
			Port:     "6543",
			User:     "postgres",
			Password: "password",
			DBName:   "trento",
		},
	}
	config, err := LoadConfig()
	suite.NoError(err)

	suite.EqualValues(expectedConfig, config)
}

func (suite *WebCmdTestSuite) TestConfigFromFlags() {
	suite.cmd.SetArgs([]string{
		"serve",
		"--host=some-host",
		"--port=1337",
		"--collector-port=1338",
		"--enable-mtls",
		"--cert=some-cert",
		"--key=some-key",
		"--ca=some-ca",
		"--db-host=some-db-host",
		"--db-port=6543",
		"--db-user=postgres",
		"--db-password=password",
		"--db-name=trento",
	})
}

func (suite *WebCmdTestSuite) TestConfigFromEnv() {
	os.Setenv("TRENTO_HOST", "some-host")
	os.Setenv("TRENTO_PORT", "1337")
	os.Setenv("TRENTO_COLLECTOR_PORT", "1338")
	os.Setenv("TRENTO_ENABLE_MTLS", "true")
	os.Setenv("TRENTO_CERT", "some-cert")
	os.Setenv("TRENTO_KEY", "some-key")
	os.Setenv("TRENTO_CA", "some-ca")
	os.Setenv("TRENTO_DB_HOST", "some-db-host")
	os.Setenv("TRENTO_DB_PORT", "6543")
	os.Setenv("TRENTO_DB_USER", "postgres")
	os.Setenv("TRENTO_DB_PASSWORD", "password")
	os.Setenv("TRENTO_DB_NAME", "trento")
}

func (suite *WebCmdTestSuite) TestConfigFromFile() {
	os.Setenv("TRENTO_CONFIG", "../../test/fixtures/config/web.yaml")
}
