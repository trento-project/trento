package runner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const (
	TrentoWebApiHost     = "TRENTO_WEB_API_HOST"
	TrentoWebApiPort     = "TRENTO_WEB_API_PORT"
	AnsibleConfigFileEnv = "ANSIBLE_CONFIG"
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

	logCommand(cmd)
	err := cmd.Run()

	if err != nil {
		log.Errorf("An error occurred while running ansible: %s", err)
		return err
	}

	log.Info("Ansible playbook execution finished successfully")

	return nil
}

func logCommand(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	go func() {
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			log.Infof(in.Text())
		}
	}()
	go func() {
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			log.Debugf(in.Text())
		}
	}()
}
