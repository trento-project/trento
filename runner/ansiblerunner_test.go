package runner

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/runner/mocks"
)

func TestLoadAraPlugins(t *testing.T) {

	a := DefaultAnsibleRunner()

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

	err := a.LoadAraPlugins()
	a.SetAraServer("127.0.0.1")

	expectedMetaRunner := &AnsibleRunner{
		Playbook: "main.yml",
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

func TestRunPlaybookSimple(t *testing.T) {

	runnerInst := &AnsibleRunner{
		Playbook: "superplay.yml",
	}

	cmd := exec.Command("echo", "stdout", "&&", ">&2", "echo", "stderr")

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute
	mockCommand.On("Execute", "ansible-playbook", "superplay.yml").Return(
		cmd,
	)

	err := runnerInst.RunPlaybook()

	assert.Equal(t, os.Environ(), cmd.Env)
	assert.NoError(t, err)

	mockCommand.AssertExpectations(t)
}

func TestRunPlaybookError(t *testing.T) {

	runnerInst := &AnsibleRunner{
		Playbook: "superplay.yml",
	}

	cmd := exec.Command("error")

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute
	mockCommand.On("Execute", "ansible-playbook", "superplay.yml").Return(
		cmd,
	)

	err := runnerInst.RunPlaybook()

	assert.Equal(t, os.Environ(), cmd.Env)
	assert.EqualError(t, err, "exec: \"error\": executable file not found in $PATH")

	mockCommand.AssertExpectations(t)
}

func TestRunPlaybookComplex(t *testing.T) {

	runnerInst := &AnsibleRunner{
		Playbook:  "superplay.yml",
		Inventory: "inventory.yml",
		Envs: map[string]string{
			"env1": "value1",
			"env2": "value2",
		},
		Check: true,
	}

	cmd := exec.Command("echo", "stdout", "&&", ">&2", "echo", "stderr")

	mockCommand := new(mocks.CustomCommand)
	customExecCommand = mockCommand.Execute
	mockCommand.On(
		"Execute", "ansible-playbook", "superplay.yml",
		"--inventory=inventory.yml", "--check").Return(
		cmd,
	)

	runnerInst.SetConfigFile("/path/myconfig.conf")

	err := runnerInst.RunPlaybook()

	assert.Contains(t, cmd.Env, "env1=value1")
	assert.Contains(t, cmd.Env, "env2=value2")
	assert.Contains(t, cmd.Env, "ANSIBLE_CONFIG=/path/myconfig.conf")
	assert.NoError(t, err)

	mockCommand.AssertExpectations(t)
}
