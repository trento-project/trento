package agent

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	consulAgent "github.com/hashicorp/consul/agent"
	configAgent "github.com/hashicorp/consul/agent/config"
	"github.com/hashicorp/go-uuid"
)

const (
	CONSUL_DNS_PORT      = 8600
	CONSUL_HTTP_PORT     = 8500
	CONSUL_HTTPS_PORT    = -1 // Disabled
	CONSUL_SERF_LAN_PORT = 8301
	CONSUL_SERF_WAN_PORT = 8302
	CONSUL_SERVER_PORT   = 8300
	CONSUL_GRPC_PORT     = 8502
)

func GetNodeName() string {
	id, err := os.Hostname()
	if err != nil {
		id, err := uuid.GenerateUUID()
		if err != nil {
			panic(err)
		}
		return id
	}

	return id
}

func getDefaultPorts() (data string) {
	return `
		ports = {
			dns = ` + strconv.Itoa(CONSUL_DNS_PORT) + `
			http = ` + strconv.Itoa(CONSUL_HTTP_PORT) + `
			https = ` + strconv.Itoa(CONSUL_HTTPS_PORT) + `
			serf_lan = ` + strconv.Itoa(CONSUL_SERF_LAN_PORT) + `
			serf_wan = ` + strconv.Itoa(CONSUL_SERF_WAN_PORT) + `
			server = ` + strconv.Itoa(CONSUL_SERVER_PORT) + `
			grpc = ` + strconv.Itoa(CONSUL_GRPC_PORT) + `
		}
	`
}

func getConfigHCL(bindAddr net.IP, srvAddr net.IP) string {
	return fmt.Sprintf(`
		bind_addr = "%s"		
		retry_join = ["%s"]		
		`, bindAddr, srvAddr,
	)
}

func NewConsulAgent(bindAddr net.IP, srvAddr net.IP) (*consulAgent.Agent, error) {
	portsConfig := getDefaultPorts()
	log.Println(portsConfig)
	d := filepath.ToSlash("./consul-agent-data")
	hclDataDir := fmt.Sprintf(`data_dir = "%s"`, d)
	consulHCLConfig := getConfigHCL(bindAddr, srvAddr)

	loader := func(source configAgent.Source) (*configAgent.RuntimeConfig, []string, error) {
		opts := configAgent.BuilderOpts{
			HCL: []string{consulHCLConfig, portsConfig, "", hclDataDir},
		}
		overrides := []configAgent.Source{}
		cfg, warnings, err := configAgent.Load(opts, source, overrides...)
		if cfg != nil {
			cfg.Telemetry.Disable = true
		}
		return cfg, warnings, err
	}
	bd, err := consulAgent.NewBaseDeps(loader, log.Writer())

	if err != nil {
		return nil, fmt.Errorf("failed to create base deps: %w", err)
	}

	return consulAgent.New(bd)
}

func (a *Agent) startConsulAgent(status chan int) error {
	a.consulAgent.Start(context.Background())
	defer a.stopConsulAgent(status)
	<-a.ctx.Done()
	return nil
}

func (a *Agent) stopConsulAgent(status chan int) {
	log.Println("Stopping consul-agent")
	a.consul.Agent().ServiceDeregister(a.cfg.InstanceName)
	a.consulAgent.Leave()
	a.consulAgent.ShutdownAgent()
	a.consulAgent.ShutdownEndpoints()
	log.Println("Consul agent stopped")
}
