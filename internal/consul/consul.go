package consul

import (
	consulApi "github.com/hashicorp/consul/api"
)

//go:generate mockgen -destination ../../test/mock_consul/consul.go github.com/SUSE/console-for-sap-applications/internal/consul Client,Catalog

type Client interface {
	Catalog() Catalog
}

type Catalog interface {
	Datacenters() ([]string, error)
	Nodes(q *consulApi.QueryOptions) ([]*consulApi.Node, *consulApi.QueryMeta, error)
}

func DefaultClient() (Client, error)  {
	w, err := consulApi.NewClient(consulApi.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &client{w}, nil
}

type client struct {
	wrapped *consulApi.Client
}

func (c *client) Catalog() Catalog {
	return c.wrapped.Catalog()
}


