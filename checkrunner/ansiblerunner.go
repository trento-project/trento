package checkrunner

import (
  "fmt"
  "os"
  "exec"
  "strings"

  log "github.com/sirupsen/logrus"
)

const (
  AnsibleCallbackPlugins = "ANSIBLE_CALLBACK_PLUGINS"
  AnsibleActionPlugins = "ANSIBLE_ACTION_PLUGINS"
)

type AnsibleRunner struct {
  Playbook  string
  Inventory string
  Envs      []string
}

func runPlaybook(playbook, inventory string, envs []string) error {
  log.Infof("Running ansible playbook %s, with inventory %s...", playbook, inventory)

  cmd := exec.Command("ansible-playbook", playbook, fmt.Sprintf("--inventory=%s", inventory))

  for _, e := range envs {
    log.Debugf("New environment variable: %s", e)
    cmd.Env = append(os.Environ(), e)
  }

  output, err := cmd.CombinedOutput()

  log.Debugf("Ansible output:\n%s:", output)

  if err != nil {
    log.Errorf("An error occurred while running ansible:", err)
    return err
  }

  log.Info("Ansible playbook execution finished successfully")

  return nil
}

func NewAnsibleRunner(playbook, inventory string) (*AnsibleRunner, error){
  r := &AnsibleRunner{}

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

  a.Envs = append(a.Envs, fmt.Sprintf("%s=%s", AnsibleCallbackPlugins, araCallbackPathStr))

	araAction := exec.Command("python3", "-m", "ara.setup.action_plugins")
	araActionPath, err := araAction.Output()
	if err != nil {
		log.Errorf("An error occurred getting the ARA actions plugin path: %s", err)
    return err
	}
	araActionPathStr := strings.TrimSpace(string(araActionPath))
  log.Debugf("ARA actions plugin found: %s", araActionPathStr)

  a.Envs = append(a.Envs, fmt.Sprintf("%s=%s", AnsibleActionPlugins, araActionPathStr))

  log.Info("ARA plugins loaded successfully")

  return nil
}

func (a *AnsibleRunner) RunPlaybook() error {
  return runPlaybook(a.Playbook, a.Inventory, a.Envs)
}
