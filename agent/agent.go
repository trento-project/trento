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
	consul        *consul.Client
	consulService *consul.AgentServiceRegistration
	ctx           context.Context
	Stop          context.CancelFunc
}

type checkFunc func() (string, error)

func New() (*Agent, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "could not read the hostname")
	}

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Consul client")
	}

	consulService := &consul.AgentServiceRegistration{
		Name: name,
		Check: &consul.AgentServiceCheck{
			TTL: "10s",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	agent := &Agent{client, consulService, ctx, cancel}
	return agent, nil
}

func (a *Agent) Start() error {
	var err error
	defer a.deregister()

	log.Println("Registering the agent service with Consul...")
	err = a.consul.Agent().ServiceRegister(a.consulService)
	if err != nil {
		return errors.Wrap(err, "could not register the agent service with Consul")
	}
	log.Println("Consul service registered.")

	log.Println("Starting Consul TTL Check loop...")
	err = a.startConsulCheckTicker(func() (string, error) {
		return "all is good", nil
	})
	if err != nil {
		return errors.Wrap(err, "could not start Consul TTL Check loop")
	}
	log.Println("Consul TTL Check loop stopped.")

	return nil
}
func (a *Agent) deregister() {
	log.Println("De-registering the agent service with Consul...")
	err := a.consul.Agent().ServiceDeregister(a.consulService.Name)
	if err != nil {
		log.Println("An error occurred while trying to deregister the agent service with Consul: ", err)
		return
	}
	log.Println("Consul service de-registered.")
}

func (a *Agent) startConsulCheckTicker(check checkFunc) error {
	a.consulCheck(check) // immediate first tick

	duration, err := time.ParseDuration(a.consulService.Check.TTL)
	if err != nil {
		return errors.Wrap(err, "could not parse TTL duration")
	}

	ticker := time.NewTicker(duration / 2)
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

func (a *Agent) consulCheck(check checkFunc) {
	log.Println("Updating Consul check TTL...")

	var err error
	var ttlStatus string

	checkOutput, checkErr := check()
	if checkErr != nil {
		checkOutput = checkErr.Error()
		ttlStatus = consul.HealthCritical
	} else {
		ttlStatus = consul.HealthPassing
	}

	err = a.consul.Agent().UpdateTTL("service:"+a.consulService.Name, checkOutput, ttlStatus)
	if err != nil {
		log.Println("An error occurred while trying to update TTL with Consul: ", err)
	}

	log.Printf("Consul check TTL updated. Status: %s.", ttlStatus)
}
