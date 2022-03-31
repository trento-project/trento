package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/agent/discovery"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal"
)

const trentoAgentCheckId = "trentoAgent"

type Agent struct {
	config          *Config
	collectorClient collector.Client
	discoveries     []discovery.Discovery
	ctx             context.Context
	ctxCancel       context.CancelFunc
}

type Config struct {
	InstanceName      string
	DiscoveriesConfig *discovery.DiscoveriesConfig
}

// NewAgent returns a new instance of Agent with the given configuration
func NewAgent(config *Config) (*Agent, error) {
	collectorClient, err := collector.NewCollectorClient(config.DiscoveriesConfig.CollectorConfig)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a collector client")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	var discoveries discovery.DiscoveryList
	discoveries, err = discoveries.AddDiscovery(discovery.NewClusterDiscovery, collectorClient, *config.DiscoveriesConfig)
	discoveries, err = discoveries.AddDiscovery(discovery.NewSAPSystemsDiscovery, collectorClient, *config.DiscoveriesConfig)
	discoveries, err = discoveries.AddDiscovery(discovery.NewCloudDiscovery, collectorClient, *config.DiscoveriesConfig)
	discoveries, err = discoveries.AddDiscovery(discovery.NewSubscriptionDiscovery, collectorClient, *config.DiscoveriesConfig)
	discoveries, err = discoveries.AddDiscovery(discovery.NewHostDiscovery, collectorClient, *config.DiscoveriesConfig)

	agent := &Agent{
		config:          config,
		collectorClient: collectorClient,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		discoveries:     discoveries,
	}
	return agent, nil
}

// Start the Agent. This will start the discovery ticker and the heartbeat ticker
func (a *Agent) Start() error {
	var wg sync.WaitGroup

	for _, d := range a.discoveries {
		wg.Add(1)
		go func(wg *sync.WaitGroup, d discovery.Discovery) {
			log.Infof("Starting %s loop...", d.GetId())
			defer wg.Done()
			a.startDiscoverTicker(d)
			log.Infof("%s discover loop stopped.", d.GetId())
		}(&wg, d)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Info("Starting heartbeat loop...")
		defer wg.Done()
		a.startHeartbeatTicker()
		log.Info("heartbeat loop stopped.")
	}(&wg)

	wg.Wait()

	return nil
}

func (a *Agent) Stop() {
	a.ctxCancel()
}

// Start a Ticker loop that will iterate over the hardcoded list of Discovery backends and execute them.
func (a *Agent) startDiscoverTicker(d discovery.Discovery) {

	tick := func() {
		result, err := d.Discover()
		if err != nil {
			result = fmt.Sprintf("Error while running discovery '%s': %s", d.GetId(), err)
			log.Errorln(result)
		}
		log.Infof("%s discovery tick output: %s", d.GetId(), result)
	}
	internal.Repeat(d.GetId(), tick, time.Duration(d.GetInterval()), a.ctx)

}

func (a *Agent) startHeartbeatTicker() {
	tick := func() {
		err := a.collectorClient.Heartbeat()
		if err != nil {
			log.Errorf("Error while sending the heartbeat to the server: %s", err)
		}
	}

	internal.Repeat("agent.heartbeat", tick, internal.HeartbeatInterval, a.ctx)
}
