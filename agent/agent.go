package agent

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul-template/manager"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/agent/discovery"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/version"
)

type Agent struct {
	cfg            Config
	discoveries    []discovery.Discovery
	consul         consul.Client
	ctx            context.Context
	ctxCancel      context.CancelFunc
	templateRunner *manager.Runner
}

type Config struct {
	InstanceName    string
	ConsulConfigDir string
	DiscoveryPeriod time.Duration
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
	client, err := consul.DefaultClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not create a Consul client")
	}

	templateRunner, err := NewTemplateRunner(cfg.ConsulConfigDir)
	if err != nil {
		return nil, errors.Wrap(err, "could not create the consul template runner")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	agent := &Agent{
		cfg:       cfg,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		consul:    client,
		discoveries: []discovery.Discovery{
			discovery.NewClusterDiscovery(client),
			discovery.NewSAPSystemsDiscovery(client),
			discovery.NewCloudDiscovery(client),
			discovery.NewSubscriptionDiscovery(client),
		},
		templateRunner: templateRunner,
	}
	return agent, nil
}

func DefaultConfig() (Config, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Config{}, errors.Wrap(err, "could not read the hostname")
	}

	return Config{
		InstanceName:    hostname,
		DiscoveryPeriod: 2 * time.Minute,
	}, nil
}

// Start the Agent which includes the registration against Consul Agent
func (a *Agent) Start() error {
	log.Println("Registering the agent service with Consul...")
	err := a.registerConsulService()
	if err != nil {
		return errors.Wrap(err, "could not register consul service")
	}
	log.Println("Consul service registered.")

	defer func() {
		log.Println("De-registering the agent service with Consul...")
		err := a.consul.Agent().ServiceDeregister(a.cfg.InstanceName)
		if err != nil {
			log.Println("An error occurred while trying to deregisterConsulService the agent service with Consul:", err)
		} else {
			log.Println("Consul service de-registered.")
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	// The Discover Loop is executing at a much slower pace than the Checker Loop
	// and will keep namespaces in Key-Value Consul store updated with specific facts
	// discovered on the node
	go func(wg *sync.WaitGroup) {
		log.Println("Starting Discover loop...")
		defer wg.Done()
		a.startDiscoverTicker()
		log.Println("Discover loop stopped.")
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Println("Starting consul-template loop...")
		defer wg.Done()
		a.startConsulTemplate()
		log.Println("consul-template loop stopped.")
	}(&wg)

	storeAgentMetadata(a.consul, version.Version)

	wg.Wait()

	return nil
}

func (a *Agent) Stop() {
	a.ctxCancel()
}

func (a *Agent) registerConsulService() error {
	var err error

	discoveryTTL := a.cfg.DiscoveryPeriod * 2
	consulService := &consulApi.AgentServiceRegistration{
		ID:   a.cfg.InstanceName,
		Name: "trento-agent",
		Tags: []string{"trento"},
		Checks: consulApi.AgentServiceChecks{
			&consulApi.AgentServiceCheck{
				CheckID: discovery.ClusterDiscoveryId,
				Name:    "HA Cluster Discovery",
				Notes:   "Collects details about the HA Cluster components running on this node",
				TTL:     discoveryTTL.String(),
				Status:  consulApi.HealthWarning,
			},
			&consulApi.AgentServiceCheck{
				CheckID: discovery.SAPDiscoveryId,
				Name:    "SAP System Discovery",
				Notes:   "Collects details about SAP System components running on this node",
				TTL:     discoveryTTL.String(),
				Status:  consulApi.HealthWarning,
			},
			&consulApi.AgentServiceCheck{
				CheckID: discovery.CloudDiscoveryId,
				Name:    "Cloud metadata discovery",
				Notes:   "Collects details about the cloud instance metadata",
				TTL:     discoveryTTL.String(),
				Status:  consulApi.HealthWarning,
			},
			&consulApi.AgentServiceCheck{
				CheckID: discovery.SubscriptionDiscoveryId,
				Name:    "Subscription discovery",
				Notes:   "Collects details about the SUSE subscription",
				TTL:     discoveryTTL.String(),
				Status:  consulApi.HealthWarning,
			},
		},
	}

	err = a.consul.Agent().ServiceRegister(consulService)
	if err != nil {
		return errors.Wrap(err, "could not register the agent service with Consul")
	}

	return nil
}

// Start a Ticker loop that will iterate over the hardcoded list of Discovery backends
// and execute them. The initial run will happen relatively quickly after Agent launch
// subsequent runs are done with a 15 minute delay. The effectiveness of the discoveries
// is reported back in the "discover_cluster" Service in consul under a TTL of 60 minutes
func (a *Agent) startDiscoverTicker() {
	tick := func() {
		for _, d := range a.discoveries {
			var status, result string
			var err error

			result, err = d.Discover()
			if err != nil {
				log.Printf("Error while running discovery '%s': %s", d.GetId(), err)
				result = err.Error()
				status = consulApi.HealthWarning
			} else {
				status = consulApi.HealthPassing
			}

			err = a.consul.Agent().UpdateTTL(d.GetId(), result, status)
			if err != nil {
				log.Println("An error occurred while trying to update TTL with Consul:", err)
			}
		}
	}

	interval := a.cfg.DiscoveryPeriod

	repeat(tick, interval, a.ctx)
}

func repeat(tick func(), interval time.Duration, ctx context.Context) {
	// run the first tick immediately
	tick()

	ticker := time.NewTicker(interval)
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

func storeAgentMetadata(client consul.Client, version string) error {
	metadata := hosts.Metadata{
		AgentVersion: version,
	}

	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
