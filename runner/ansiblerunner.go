// Find ARA ansible resources information at:
// https://ara.readthedocs.io/en/latest/
// https://github.com/ansible-community/ara
// https://ara.readthedocs.io/en/latest/index.html

package runner

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	TrentoWebApiHost       = "TRENTO_WEB_API_HOST"
	TrentoWebApiPort       = "TRENTO_WEB_API_PORT"
	AnsibleConfigFileEnv   = "ANSIBLE_CONFIG"
	AnsibleCallbackPlugins = "ANSIBLE_CALLBACK_PLUGINS"
	AnsibleActionPlugins   = "ANSIBLE_ACTION_PLUGINS"
	AraApiClient           = "ARA_API_CLIENT"
	AraApiServer           = "ARA_API_SERVER"
)

//go:generate mockery --name=CustomCommand

type CustomCommand func(name string, arg ...string) *exec.Cmd

var customExecCommand CustomCommand = exec.Command

type AnsibleRunner struct {
	Playbook  string
	Inventory string
	Envs      map[string]string
	Check     bool
}

func DefaultAnsibleRunner() *AnsibleRunner {
	return &AnsibleRunner{
		Playbook: "main.yml",
		Envs:     make(map[string]string),
		Check:    false,
	}
}

func DefaultAnsibleRunnerWithAra() (*AnsibleRunner, error) {
	a := DefaultAnsibleRunner()
	if err := a.LoadAraPlugins(); err != nil {
		return a, err
	}

	return a, nil
}

func (a *AnsibleRunner) setEnv(name, value string) {
	a.Envs[name] = value
}

func (a *AnsibleRunner) SetPlaybook(playbook string) error {
	if _, err := os.Stat(playbook); os.IsNotExist(err) {
		log.Errorf("Playbook file %s does not exist", playbook)
		return err
	}

	a.Playbook = playbook
	return nil
}

func (a *AnsibleRunner) SetInventory(inventory string) error {
	if _, err := os.Stat(inventory); os.IsNotExist(err) {
		log.Errorf("Inventory file %s does not exist", inventory)
		return err
	}

	a.Inventory = inventory
	return nil
}

func (a *AnsibleRunner) SetConfigFile(confFile string) {
	a.setEnv(AnsibleConfigFileEnv, confFile)
}

func (a *AnsibleRunner) SetTrentoApiData(host string, port int) {
	a.setEnv(TrentoWebApiHost, host)
	a.setEnv(TrentoWebApiPort, fmt.Sprintf("%d", port))
}

// ARA_API_CLIENT is always set to "http" to ensure the usage of the REST API
// "offline" mode could be used, but it would only work if the ansible runner is
// running in the same host, and it doesn't provide much value
func (a *AnsibleRunner) SetAraServer(host string) {
	a.setEnv(AraApiClient, "http")
	a.setEnv(AraApiServer, host)
}

func (a *AnsibleRunner) LoadAraPlugins() error {
	log.Info("Loading ARA plugins...")

	araCallback := customExecCommand("python3", "-m", "ara.setup.callback_plugins")
	araCallbackPath, err := araCallback.Output()
	if err != nil {
		log.Errorf("An error occurred getting the ARA callback plugin path: %s", err)
		return err
	}
	araCallbackPathStr := strings.TrimSpace(string(araCallbackPath))
	log.Debugf("ARA callback plugin found: %s", araCallbackPathStr)

	a.setEnv(AnsibleCallbackPlugins, araCallbackPathStr)

	araAction := customExecCommand("python3", "-m", "ara.setup.action_plugins")
	araActionPath, err := araAction.Output()
	if err != nil {
		log.Errorf("An error occurred getting the ARA actions plugin path: %s", err)
		return err
	}
	araActionPathStr := strings.TrimSpace(string(araActionPath))
	log.Debugf("ARA actions plugin found: %s", araActionPathStr)

	a.setEnv(AnsibleActionPlugins, araActionPathStr)

	log.Info("ARA plugins loaded successfully")

	return nil
}

func (a *AnsibleRunner) IsAraServerUp() bool {
	server, ok := a.Envs[AraApiServer]
	if !ok {
		log.Warn("ARA server usage not configured")
		return false
	}

	host := fmt.Sprintf("%s/api/", server)
	log.Debugf("Looking for the ARA server at: %s", host)

	resp, err := http.Get(host)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Debugf("Error requesting ARA server api: %s", err)
		return false
	}

	log.Debugf("ARA server response code: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (a *AnsibleRunner) RunPlaybook() error {
	var cmdItems []string

	log.Infof("Ansible playbook %s", a.Playbook)
	cmdItems = append(cmdItems, a.Playbook)

	if a.Inventory != "" {
		log.Infof("Inventory %s", a.Inventory)
		cmdItems = append(cmdItems, fmt.Sprintf("--inventory=%s", a.Inventory))
	}

	if a.Check {
		log.Info("Running in check mode")
		cmdItems = append(cmdItems, "--check")
	}

	cmd := customExecCommand("ansible-playbook", cmdItems...)
	cmd.Env = os.Environ()
	for key, value := range a.Envs {
		newEnv := fmt.Sprintf("%s=%s", key, value)
		log.Debugf("New environment variable: %s", newEnv)
		cmd.Env = append(cmd.Env, newEnv)
	}

	output, err := cmd.CombinedOutput()
	os.Stdout.Write(output)
	log.Infof("Ansible output:\n%s:", output)

	if err != nil {
		log.Errorf("An error occurred while running ansible: %s", err)
		return err
	}

	log.Info("Ansible playbook execution finished successfully")

	return nil
}
