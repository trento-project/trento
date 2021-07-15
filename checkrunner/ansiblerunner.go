// Find ARA ansible resources information at:
// https://ara.readthedocs.io/en/latest/
// https://github.com/ansible-community/ara
// https://ara.readthedocs.io/en/latest/index.html

package checkrunner

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	AnsibleCallbackPlugins = "ANSIBLE_CALLBACK_PLUGINS"
	AnsibleActionPlugins   = "ANSIBLE_ACTION_PLUGINS"
	AraApiClient           = "ARA_API_CLIENT"
	AraApiServer           = "ARA_API_SERVER"
)

type AnsibleRunner struct {
	Playbook  string
	Inventory string
	Envs      map[string]string
}

func runPlaybook(playbook, inventory string, envs map[string]string) error {
	log.Infof("Running ansible playbook %s, with inventory %s...", playbook, inventory)

	cmd := exec.Command("ansible-playbook", playbook, fmt.Sprintf("--inventory=%s", inventory))

	cmd.Env = os.Environ()
	for key, value := range envs {
		newEnv := fmt.Sprintf("%s=%s", key, value)
		log.Debugf("New environment variable: %s", newEnv)
		cmd.Env = append(cmd.Env, newEnv)
	}

	output, err := cmd.CombinedOutput()

	log.Debugf("Ansible output:\n%s:", output)

	if err != nil {
		log.Errorf("An error occurred while running ansible: %s", err)
		return err
	}

	log.Info("Ansible playbook execution finished successfully")

	return nil
}

func NewAnsibleRunner(playbook, inventory string) (*AnsibleRunner, error) {
	r := &AnsibleRunner{Envs: make(map[string]string)}

	if _, err := os.Stat(playbook); os.IsNotExist(err) {
		log.Errorf("Playbook file %s does not exist", playbook)
		return r, err
	}

	r.Playbook = playbook

	if _, err := os.Stat(inventory); os.IsNotExist(err) {
		log.Errorf("Inventory file %s does not exist", inventory)
		return r, err
	}

	r.Inventory = inventory

	return r, nil
}

func (a *AnsibleRunner) setEnv(name, value string) {
	a.Envs[name] = value
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

	araCallback := exec.Command("python3", "-m", "ara.setup.callback_plugins")
	araCallbackPath, err := araCallback.Output()
	if err != nil {
		log.Errorf("An error occurred getting the ARA callback plugin path: %s", err)
		return err
	}
	araCallbackPathStr := strings.TrimSpace(string(araCallbackPath))
	log.Debugf("ARA callback plugin found: %s", araCallbackPathStr)

	a.setEnv(AnsibleCallbackPlugins, araCallbackPathStr)

	araAction := exec.Command("python3", "-m", "ara.setup.action_plugins")
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

func (a *AnsibleRunner) isAraServerUp() bool {
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
	return runPlaybook(a.Playbook, a.Inventory, a.Envs)
}
