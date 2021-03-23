package consul

import (
	consulApi "github.com/hashicorp/consul/api"
)

type Client interface {
	Catalog() Catalog
}

type Catalog interface {
	Datacenters() ([]string, error)
	Nodes(q *consulApi.QueryOptions) ([]*consulApi.Node, *consulApi.QueryMeta, error)
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
