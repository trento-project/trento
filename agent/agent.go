package agent

import (
	"context"
	"log"
	"os"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Agent struct {
	cfg       Config
	check     Check
	consul    *consul.Client
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type Config struct {
	InstanceName    string
	DefinitionsPath string
	TTL             time.Duration
}

func New() (*Agent, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, errors.Wrap(err, "could not create the agent configuration")
	}

	return NewWithConfig(config)
}

// returns a new instance of Agent with the given configuration
func NewWithConfig(cfg Config) (*Agent, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Consul client")
	}

	checker, err := NewCheck(cfg.DefinitionsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Checker instance")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	agent := &Agent{
		cfg:       cfg,
		check:     checker,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		consul:    client,
	}
	return agent, nil
}

func DefaultConfig() (Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Config{}, errors.Wrap(err, "could not read the hostname")
	}

	return Config{
		InstanceName: hostname,
		TTL:          10 * time.Second,
	}, nil
}

func (a *Agent) Start() error {
	var err error
	defer a.deregister()

	log.Println("Registering the agent service with Consul...")

	consulService := &consul.AgentServiceRegistration{
		Name: a.cfg.InstanceName,
		Tags: []string{"console-agent"},
		Check: &consul.AgentServiceCheck{
			Name:   "HA config checks",
			Notes:  "Checks whether or not the HA configuration is compliant with the provided best practices",
			TTL:    a.cfg.TTL.String(),
			Status: consul.HealthWarning,
		},
	}

	err = a.consul.Agent().ServiceRegister(consulService)
	if err != nil {
		return errors.Wrap(err, "could not register the agent service with Consul")
	}
	log.Println("Consul service registered.")

	log.Println("Starting Consul TTL Check loop...")
	err = a.startConsulCheckTicker(a.check)
	if err != nil {
		return errors.Wrap(err, "could not start Consul TTL Check loop")
	}
	log.Println("Consul TTL Check loop stopped.")

	return nil
}

func (a *Agent) Stop() {
	a.ctxCancel()
}

func (a *Agent) deregister() {
	log.Println("De-registering the agent service with Consul...")
	err := a.consul.Agent().ServiceDeregister(a.cfg.InstanceName)
	if err != nil {
		log.Println("An error occurred while trying to deregister the agent service with Consul:", err)
		return
	}
	log.Println("Consul service de-registered.")
}

func (a *Agent) startConsulCheckTicker(check Check) error {
	a.consulCheck(check) // immediate first tick

	ticker := time.NewTicker(a.cfg.TTL / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.consulCheck(check)
		case <-a.ctx.Done():
			return nil
		}
	}
}

func (a *Agent) consulCheck(check Check) {
	log.Println("Updating Consul check TTL...")

	var err error

	result, err := check()
	if err != nil {
		log.Println("An error occurred while running health checks:", err)
		return
	}

	err = a.consul.Agent().UpdateTTL("service:"+a.cfg.InstanceName, result.String(), result.Status)
	if err != nil {
		log.Println("An error occurred while trying to update TTL with Consul:", err)
		return
	}

	log.Printf("Consul check TTL updated. Status: %s.", result.Status)
}
