package agent

import (
	"context"
	"log"
	"os"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Agent struct {
	consul        *consulApi.Client
	consulService *consulApi.AgentServiceRegistration
	ctx           context.Context
	cancel        context.CancelFunc
}

type checkFunc func() (string, error)

func New() (*Agent, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "could not read the hostname")
	}

	consul, err := consulApi.NewClient(consulApi.DefaultConfig())
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Consul client")
	}

	consulService := &consulApi.AgentServiceRegistration{
		Name: name,
		Check: &consulApi.AgentServiceCheck{
			TTL: "10s",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	agent := &Agent{consul, consulService, ctx, cancel}
	return agent, nil
}

func (a *Agent) Start() error {
	log.Println("Starting the Console Agent...")

	var err error

	log.Println("Registering the agent service with Consul...")
	err = a.consul.Agent().ServiceRegister(a.consulService)
	if err != nil {
		return errors.Wrap(err, "could not register the agent service with Consul")
	}
	log.Println("Consul service registered.")

	a.startConsulCheckTicker(func() (string, error) {
		return "all is good", nil
	})

	return nil
}

func (a *Agent) Stop() {
	a.cancel()

	log.Println("De-registering the agent service with Consul...")
	err := a.consul.Agent().ServiceDeregister(a.consulService.Name)
	if err != nil {
		log.Println("An error occurred while trying to deregister the agent service with Consul: ", err)
		return
	}
	log.Println("Consul service de-registered.")
}

func (a *Agent) startConsulCheckTicker(check checkFunc) {
	a.consulCheck(check) // immediate first tick

	duration, _ := time.ParseDuration(a.consulService.Check.TTL)
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.consulCheck(check)
		case <-a.ctx.Done():
			return
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
		ttlStatus = consulApi.HealthCritical
	} else {
		ttlStatus = consulApi.HealthPassing
	}

	err = a.consul.Agent().UpdateTTL("service:"+a.consulService.Name, checkOutput, ttlStatus)
	if err != nil {
		log.Println("An error occurred while trying to update TTL with Consul: ", err)
	}

	log.Printf("Consul check TTL updated. Status: %s", ttlStatus)
}
