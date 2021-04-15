package consul

import (
	consulApi "github.com/hashicorp/consul/api"
)

//go:generate mockery --all

type Client interface {
	Catalog() Catalog
	Health() Health
}

type Catalog interface {
	Datacenters() ([]string, error)
	Node(node string, q *consulApi.QueryOptions) (*consulApi.CatalogNode, *consulApi.QueryMeta, error)
	Nodes(q *consulApi.QueryOptions) ([]*consulApi.Node, *consulApi.QueryMeta, error)
}

type Health interface {
	Node(node string, q *consulApi.QueryOptions) (consulApi.HealthChecks, *consulApi.QueryMeta, error)
	Service(service, tag string, passingOnly bool, q *consulApi.QueryOptions) ([]*consulApi.ServiceEntry, *consulApi.QueryMeta, error)
}

func DefaultClient() (Client, error) {
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

func (c *client) Health() Health {
	return c.wrapped.Health()
}
