package runner

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/runner/mocks"
)

// TODO: This test could be improved to check the definitve ansible files structure
// once we have something fixed
func TestCreateAnsibleFiles(t *testing.T) {
	tmpDir, _ := ioutil.TempDir(os.TempDir(), "trentotest")
	err := createAnsibleFiles(tmpDir)

	assert.DirExists(t, path.Join(tmpDir, "ansible"))
	assert.NoError(t, err)

	os.RemoveAll(tmpDir)
}

func TestNewAnsibleMetaRunner(t *testing.T) {

	cfg := &Config{
		AnsibleFolder: "../test/ansible_test",
		AraServer:     "127.0.0.1",
	}

	cmdCallback := exec.Command("echo", "callback")
	cmdAction := exec.Command("echo", "action")

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute
	mockCommand.On("Execute", "python3", "-m", "ara.setup.callback_plugins").Return(
		cmdCallback,
	)
	mockCommand.On("Execute", "python3", "-m", "ara.setup.action_plugins").Return(
		cmdAction,
	)

	a, err := NewAnsibleMetaRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook: "../test/ansible_test/ansible/meta.yml",
		Envs: map[string]string{
			"ANSIBLE_CALLBACK_PLUGINS": "callback",
			"ANSIBLE_ACTION_PLUGINS":   "action",
			"ARA_API_CLIENT":           "http",
			"ARA_API_SERVER":           "127.0.0.1",
		},
		Check: false,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMetaRunner, a)

	mockCommand.AssertExpectations(t)
}

func TestNewAnsibleCheckRunner(t *testing.T) {

	cfg := &Config{
		AnsibleFolder: "../test/ansible_test",
		AraServer:     "127.0.0.1",
	}

	cmdCallback := exec.Command("echo", "callback")
	cmdAction := exec.Command("echo", "action")

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute
	mockCommand.On("Execute", "python3", "-m", "ara.setup.callback_plugins").Return(
		cmdCallback,
	)
	mockCommand.On("Execute", "python3", "-m", "ara.setup.action_plugins").Return(
		cmdAction,
	)

	a, err := NewAnsibleCheckRunner(cfg)

	expectedMetaRunner := &AnsibleRunner{
		Playbook:  "../test/ansible_test/ansible/check.yml",
		Inventory: "../test/ansible_test/ansible_hosts",
		Envs: map[string]string{
			"ANSIBLE_CALLBACK_PLUGINS": "callback",
			"ANSIBLE_ACTION_PLUGINS":   "action",
			"ARA_API_CLIENT":           "http",
			"ARA_API_SERVER":           "127.0.0.1",
		},
		Check: true,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMetaRunner, a)

	mockCommand.AssertExpectations(t)
}
