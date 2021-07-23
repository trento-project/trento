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

	"github.com/hashicorp/consul-template/manager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:embed ansible
var ansibleFS embed.FS

const (
	AnsibleMain = "ansible/main.yaml"
)

type Runner struct {
	cfg            Config
	ctx            context.Context
	ctxCancel      context.CancelFunc
	templateRunner *manager.Runner
}

type Config struct {
	AraServer              string
	ConsulAddr             string
	Interval               time.Duration
	AnsibleFolder          string
	ConsulTemplateLogLevel string
}

func NewWithConfig(cfg Config) (*Runner, error) {
	templateRunner, err := NewTemplateRunner(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not create the consul template runner")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	runner := &Runner{
		cfg:            cfg,
		ctx:            ctx,
		ctxCancel:      ctxCancel,
		templateRunner: templateRunner,
	}

	return runner, nil
}

func DefaultConfig() (Config, error) {
	return Config{
		AraServer:              "http://127.0.0.1:8000",
		ConsulAddr:             "127.0.0.1:8500",
		Interval:               5 * time.Minute,
		AnsibleFolder:          "/usr/etc/trento",
		ConsulTemplateLogLevel: "info",
	}, nil
}

func (c *Runner) Start() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Println("Starting the runner loop...")
		defer wg.Done()
		c.startRunnerTicker()
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

	err = os.MkdirAll(ansibleFolder, 0644)
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
			os.Mkdir(path.Join(folder, fileName), 0644)
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

func (c *Runner) startRunnerTicker() {
	err := createAnsibleFiles(c.cfg.AnsibleFolder)
	if err != nil {
		return
	}

	c.startConsulTemplate()

	ansibleRunner, err := NewAnsibleRunner(
		path.Join(c.cfg.AnsibleFolder, AnsibleMain),
		path.Join(c.cfg.AnsibleFolder, ansibleHostFile))
	if err != nil {
		return
	}

	err = ansibleRunner.LoadAraPlugins()
	if err != nil {
		return
	}

	ansibleRunner.SetAraServer(c.cfg.AraServer)

	tick := func() {
		// As consul-template is executed as run-once, we need to create the runner everytime
		tmpRunner, err := NewTemplateRunner(&c.cfg)
		if err != nil {
			return
		}
		c.templateRunner = tmpRunner
		c.startConsulTemplate()

		if !ansibleRunner.isAraServerUp() {
			log.Error("ARA server not found. Skipping ansible execution as the data is not recorded")
			return
		}
		ansibleRunner.RunPlaybook()
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
