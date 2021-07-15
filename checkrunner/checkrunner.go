package checkrunner

import (
	"context"
	"embed"
	//"fmt"
	//"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

//go:embed ansible
var ansibleFS embed.FS

type CheckRunner struct {
	cfg       Config
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type Config struct {
	AraServer string
	Interval  time.Duration
}

func NewWithConfig(cfg Config) (*CheckRunner, error) {
	ctx, ctxCancel := context.WithCancel(context.Background())

	runner := &CheckRunner{
		cfg:       cfg,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	return runner, nil
}

func DefaultConfig() (Config, error) {
	return Config{
		AraServer: "http://127.0.0.1:8000",
		Interval:  5 * time.Minute,
	}, nil
}

func (c *CheckRunner) Start() error {
	var wg sync.WaitGroup

	wg.Add(1)

	//createTempAnsible()

	go func(wg *sync.WaitGroup) {
		log.Println("Starting the check runner loop...")
		defer wg.Done()
		c.startCheckRunnerTicker()
		log.Println("Check runner loop stopped.")
	}(&wg)

	wg.Wait()

	return nil
}

func (c *CheckRunner) Stop() {
	c.ctxCancel()
}

/*
func createTempAnsible() error {
	err := os.RemoveAll("consul.d/ansible")
	if err != nil {
		log.Print(err)
		return err
	}

	err = os.Mkdir("consul.d/ansible", 0644)
	if err != nil {
		log.Print(err)
		return err
	}

	err = fs.WalkDir(ansibleFS, "ansible", func(fileName string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dir.IsDir() {
			content, err := ansibleFS.ReadFile(fileName)
			if err != nil {
				log.Printf("Error reading file %s", fileName)
				return err
			}
			f, err := os.Create(path.Join("consul.d", fileName))
			if err != nil {
				log.Printf("Error creating file %s", fileName)
				return err
			}
			fmt.Fprintf(f, "%s", content)
		} else {
			os.Mkdir(path.Join("consul.d", fileName), 0644)
		}
		return nil
	})

	return nil
}
*/

func (c *CheckRunner) startCheckRunnerTicker() {
	ansibleRunner, err := NewAnsibleRunner(
		"/srv/trento/consul.d/ansible/main.yaml", "/srv/trento/consul.d/ansible_hosts")
	if err != nil {
		return
	}

	err = ansibleRunner.LoadAraPlugins()
	if err != nil {
		return
	}

	ansibleRunner.SetAraHost(c.cfg.AraServer)

	tick := func() {
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
		case <-ctx.Done():
			return
		}
	}
}
