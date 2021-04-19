package agent

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
)

type ConsulTemplateConfig struct {
	Source      string
	Destination string
}

func StartConsulTemplate(a *Agent, c *ConsulTemplateConfig) (*manager.Runner, error) {
	config := consultemplateconfig.DefaultConfig()
	*config.Templates = append(
		*config.Templates,
		&consultemplateconfig.TemplateConfig{
			Source:      &c.Source,
			Destination: &c.Destination,
		},
	)

	runner, err := manager.NewRunner(config, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not start consul-template")
	}
	go runner.Start()

	// Create new thread method to check for consul-template events. If the template is rendered
	// reload the configuration
	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case quit := <-signals:
				log.Printf("Exiting from consul-template event listener loop: %s", quit)
				break
			case <-runner.TemplateRenderedCh():
				log.Printf("Template rendered. Reloading agent configuration...")
				err := a.consul.Agent().Reload()
				if err != nil {
					log.Printf("Error reloading agent meta-data information: %s", err)
				} else {
					log.Print("Agent meta-data correctly reloaded")
				}
			}
		}
	}()

	return runner, nil
}

func StopConsulTemplate(runner *manager.Runner) {
	log.Println("Stopping consul-template")
	runner.StopImmediately()
	log.Println("Stopped consul-template")
}
