package agent

import (
	"context"
	"fmt"
	"strings"
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
	InstanceName    string
	SSHAddress      string
	DiscoveryPeriod time.Duration
	CollectorConfig *collector.Config
}

// NewAgent returns a new instance of Agent with the given configuration
func NewAgent(config *Config) (*Agent, error) {
	collectorClient, err := collector.NewCollectorClient(config.CollectorConfig)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a collector client")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	agent := &Agent{
		config:          config,
		collectorClient: collectorClient,
		ctx:             ctx,
		ctxCancel:       ctxCancel,
		discoveries: []discovery.Discovery{
			discovery.NewClusterDiscovery(collectorClient),
			discovery.NewSAPSystemsDiscovery(collectorClient),
			discovery.NewCloudDiscovery(collectorClient),
			discovery.NewSubscriptionDiscovery(collectorClient),
			discovery.NewHostDiscovery(config.SSHAddress, collectorClient),
		},
	}
	return agent, nil
}

// Start the Agent. This will start the discovery ticker and the heartbeat ticker
func (a *Agent) Start() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		log.Info("Starting Discover loop...")
		defer wg.Done()
		a.startDiscoverTicker()
		log.Info("Discover loop stopped.")
	}(&wg)

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
func (a *Agent) startDiscoverTicker() {
	tick := func() {
		var output []string
		for _, d := range a.discoveries {
			result, err := d.Discover()
			if err != nil {
				result = fmt.Sprintf("Error while running discovery '%s': %s", d.GetId(), err)

				log.Errorln(result)
			}
			output = append(output, result)
		}
		log.Infof("Discovery tick output", strings.Join(output, "\n\n"))
	}

	interval := a.config.DiscoveryPeriod
	internal.Repeat("agent.discovery", tick, interval, a.ctx)
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
