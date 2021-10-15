package runner

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/trento-project/trento/internal/consul"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:embed ansible
var ansibleFS embed.FS

const (
	AnsibleMain       = "ansible/check.yml"
	AnsibleMeta       = "ansible/meta.yml"
	AnsibleConfigFile = "ansible/ansible.cfg"
	AnsibleHostFile   = "ansible/ansible_hosts"
)

type Runner struct {
	cfg       Config
	ctx       context.Context
	ctxCancel context.CancelFunc
	consul    consul.Client
}

type Config struct {
	AraServer     string
	Interval      time.Duration
	AnsibleFolder string
}

func NewWithConfig(cfg Config) (*Runner, error) {
	client, err := consul.DefaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create the consul client connection")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	runner := &Runner{
		cfg:       cfg,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		consul:    client,
	}

	return runner, nil
}

func DefaultConfig() (Config, error) {
	return Config{
		AraServer:              "http://127.0.0.1:8000",
		Interval:               5 * time.Minute,
		AnsibleFolder:          "/tmp/trento",
	}, nil
}

func (c *Runner) Start() error {
	var wg sync.WaitGroup

	if err := createAnsibleFiles(c.cfg.AnsibleFolder); err != nil {
		return err
	}

	metaRunner, err := NewAnsibleMetaRunner(&c.cfg)
	if err != nil {
		return err
	}

	if !metaRunner.IsAraServerUp() {
		return fmt.Errorf("ARA server not available")
	}

	if err = metaRunner.RunPlaybook(); err != nil {
		return err
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Println("Starting the runner loop...")
		defer wg.Done()
		c.startCheckRunnerTicker()
		log.Println("Runner loop stopped.")
	}(&wg)

	wg.Wait()

	return nil
}

func (c *Runner) Stop() {
	c.ctxCancel()
}

func createAnsibleFiles(folder string) error {
	log.Infof("Creating the ansible file structure in %s", folder)
	// Clean the folder if it stores old files
	ansibleFolder := path.Join(folder, "ansible")
	err := os.RemoveAll(ansibleFolder)
	if err != nil {
		log.Error(err)
		return err
	}

	err = os.MkdirAll(ansibleFolder, 0755)
	if err != nil {
		log.Error(err)
		return err
	}

	// Create the ansible file structure from the FS
	err = fs.WalkDir(ansibleFS, "ansible", func(fileName string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			content, err := ansibleFS.ReadFile(fileName)
			if err != nil {
				log.Errorf("Error reading file %s", fileName)
				return err
			}
			f, err := os.Create(path.Join(folder, fileName))
			if err != nil {
				log.Errorf("Error creating file %s", fileName)
				return err
			}
			fmt.Fprintf(f, "%s", content)
		} else {
			os.Mkdir(path.Join(folder, fileName), 0755)
		}
		return nil
	})

	if err != nil {
		log.Errorf("An error ocurred during the ansible file structure creation: %s", err)
		return err
	}

	log.Info("Ansible file structure successfully created")

	return nil
}

func NewAnsibleMetaRunner(cfg *Config) (*AnsibleRunner, error) {
	playbookPath := path.Join(cfg.AnsibleFolder, AnsibleMeta)
	ansibleRunner, err := DefaultAnsibleRunnerWithAra()
	if err != nil {
		return ansibleRunner, err
	}

	if err = ansibleRunner.SetPlaybook(playbookPath); err != nil {
		return ansibleRunner, err
	}

	configFile := path.Join(cfg.AnsibleFolder, AnsibleConfigFile)
	ansibleRunner.SetConfigFile(configFile)
	ansibleRunner.SetAraServer(cfg.AraServer)

	return ansibleRunner, err
}

func NewAnsibleCheckRunner(cfg *Config) (*AnsibleRunner, error) {
	playbookPath := path.Join(cfg.AnsibleFolder, AnsibleMain)

	ansibleRunner, err := DefaultAnsibleRunnerWithAra()
	if err != nil {
		return ansibleRunner, err
	}

	if err = ansibleRunner.SetPlaybook(playbookPath); err != nil {
		return ansibleRunner, err
	}

	ansibleRunner.Check = true
	configFile := path.Join(cfg.AnsibleFolder, AnsibleConfigFile)
	ansibleRunner.SetConfigFile(configFile)
	ansibleRunner.SetAraServer(cfg.AraServer)

	return ansibleRunner, nil
}

func (c *Runner) startCheckRunnerTicker() {
	checkRunner, err := NewAnsibleCheckRunner(&c.cfg)
	if err != nil {
		return
	}

	tick := func() {
		content, err := NewClusterInventoryContent(c.consul)
		if err != nil {
			log.Errorf("Error creating the ansible inventory content")
			return
		}

		inventoryFile := path.Join(c.cfg.AnsibleFolder, AnsibleHostFile)
		err = CreateInventory(inventoryFile, content)
		if err != nil {
			log.Errorf("Error creating the ansible inventory file")
			return
		}

		if err = checkRunner.SetInventory(inventoryFile); err != nil {
			return
		}

		if !checkRunner.IsAraServerUp() {
			log.Error("ARA server not found. Skipping ansible execution as the data is not recorded")
			return
		}
		checkRunner.RunPlaybook()
	}

	interval := c.cfg.Interval

	repeat(tick, interval, c.ctx)
}

func repeat(tick func(), interval time.Duration, ctx context.Context) {
	// run the first tick immediately
	tick()

	ticker := time.NewTicker(interval)
	log.Debugf("Next execution in %s", interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tick()
			log.Debugf("Next execution in %s", interval)
		case <-ctx.Done():
			return
		}
	}
}
