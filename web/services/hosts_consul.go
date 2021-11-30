package services

import (
	"fmt"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

//go:generate mockery --name=HostsConsulService --inpackage  --filename=hosts_consul_mock.go

type HostsConsulService interface {
	GetHostMetadata(host string) (map[string]string, error)
	GetHostsBySystemId(id string) (hosts.HostList, error)
}

type hostsConsulService struct {
	consul consul.Client
}

func NewHostsConsulService(client consul.Client) HostsConsulService {
	return &hostsConsulService{consul: client}
}

func (h *hostsConsulService) GetHostMetadata(host string) (map[string]string, error) {
	hostList, err := hosts.Load(h.consul, fmt.Sprintf("Node == %s", host), nil)
	if err != nil {
		return nil, err
	}

	if len(hostList) == 0 {
		return nil, fmt.Errorf("host with name %s not found", host)
	}

	return hostList[0].TrentoMeta(), nil
}

func (h *hostsConsulService) GetHostsBySystemId(id string) (hosts.HostList, error) {
	hostList, err := hosts.Load(h.consul, fmt.Sprintf("Meta[\"trento-sap-systems-id\"] contains \"%s\"", id), nil)
	if err != nil {
		return nil, err
	}

	return hostList, nil
}
