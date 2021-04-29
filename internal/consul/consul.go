package consul

import (
	consulApi "github.com/hashicorp/consul/api"
)

//go:generate mockery --all

type Client interface {
	Catalog() Catalog
	Health() Health
	KV() KV
}

type Catalog interface {
	Datacenters() ([]string, error)
	Node(node string, q *consulApi.QueryOptions) (*consulApi.CatalogNode, *consulApi.QueryMeta, error)
	Nodes(q *consulApi.QueryOptions) ([]*consulApi.Node, *consulApi.QueryMeta, error)
}

type KV interface {
	Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error)
	List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error)
	Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error)
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

func (c *client) KV() KV {
	return c.wrapped.KV()
}

func (c *client) Health() Health {
	return c.wrapped.Health()
}
