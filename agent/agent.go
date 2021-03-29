package agent

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Agent struct {
	cfg              Config
	check            Checker
	consulResultChan chan CheckResult
	wsResultChan     chan CheckResult
	webService       *webService
	consul           *consul.Client
	ctx              context.Context
	ctxCancel        context.CancelFunc
}

type Config struct {
	WebHost         string
	WebPort         int
	ServiceName     string
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

	checker, err := NewChecker(cfg.DefinitionsPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Checker instance")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	wsResultChan := make(chan CheckResult)

	agent := &Agent{
		cfg:          cfg,
		check:        checker,
		ctx:          ctx,
		ctxCancel:    ctxCancel,
		consul:       client,
		webService:   newWebService(wsResultChan),
		wsResultChan: wsResultChan,
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
	err := a.registerConsulService()
	if err != nil {
		return errors.Wrap(err, "could not register consul service")
	}
	defer a.deregisterConsulService()

	var wg sync.WaitGroup
	wg.Add(2)
	errs := make(chan error, 2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		errs <- a.startCheckTicker()
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		errs <- a.webService.Start(a.cfg.WebHost, a.cfg.WebPort, a.ctx)
	}(&wg)

	// As soon as all the goroutines in the Waitgroup are done, close the channel where the errors are sent
	go func() {
		wg.Wait()
		close(errs)
	}()

	// Scroll the errors channel and return as soon as one of the goroutines fails.
	// This will block until the errors channel is closed.
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Agent) Stop() {
	close(a.wsResultChan)
	a.ctxCancel()
}

func (a *Agent) registerConsulService() error {
	var err error

	log.Println("Registering the agent service with Consul...")

	consulService := &consul.AgentServiceRegistration{
		ID:   a.cfg.InstanceName,
		Name: a.cfg.ServiceName,
		Tags: []string{"console-agent"},
		Checks: consul.AgentServiceChecks{
			&consul.AgentServiceCheck{
				CheckID:  "hana_tcp_check",
				Name:     "HANA TCP check",
				TCP:      "localhost:50013",
				Interval: "10s",
			},
			{
				CheckID: "ha_checks",
				Name:    "HA config checks",
				Notes:   "Checks whether or not the HA configuration is compliant with the provided best practices",
				TTL:     a.cfg.TTL.String(),
				Status:  consul.HealthWarning,
			},
		},
	}

	err = a.consul.Agent().ServiceRegister(consulService)
	if err != nil {
		return errors.Wrap(err, "could not register the agent service with Consul")
	}
	log.Println("Consul service registered.")

	return nil
}

func (a *Agent) deregisterConsulService() {
	log.Println("De-registering the agent service with Consul...")
	err := a.consul.Agent().ServiceDeregister(a.cfg.InstanceName)
	if err != nil {
		log.Println("An error occurred while trying to deregisterConsulService the agent service with Consul:", err)
		return
	}
	log.Println("Consul service de-registered.")
}

func (a *Agent) startCheckTicker() error {
	log.Println("Starting Check loop...")
	defer log.Println("Check loop stopped.")

	ticker := time.NewTicker(a.cfg.TTL / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			result, err := a.check()
			if err != nil {
				log.Println("An error occurred while running health checks:", err)
			}
			go func() {
				select {
				case <-a.ctx.Done():
					return
				default:
					a.wsResultChan <- result
				}
			}()
			a.updateConsulCheck(result)
		case <-a.ctx.Done():
			return nil
		}
	}
}

func (a *Agent) updateConsulCheck(result CheckResult) {
	log.Println("Updating Consul check TTL...")

	var err error
	var status string

	summary := result.Summary()

	switch true {
	case summary.Fail > 0:
		status = consul.HealthCritical
	case summary.Warn > 0:
		status = consul.HealthWarning
	default:
		status = consul.HealthPassing
	}

	err = a.consul.Agent().UpdateTTL("service:"+a.cfg.InstanceName, result.String(), status)
	if err != nil {
		log.Println("An error occurred while trying to update TTL with Consul:", err)
		return
	}

	log.Printf("Consul check TTL updated. Status: %s.", status)
}
