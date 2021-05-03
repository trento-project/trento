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
	Maps(prefix, offset string) (map[string]interface{}, error)
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

func (c *client) KV() KV {
	return &kv{c.wrapped.KV()}
}

type kv struct {
	kv *consulApi.KV
}

func (k *kv) Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error) {
	return k.kv.Get(key, q)
}

func (k *kv) Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error) {
	return k.kv.Keys(prefix, separator, q)
}


func (k *kv) List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
	return k.kv.List(prefix, q)
}

func (k *kv) Maps(prefix, offset string) (map[string]interface{}, error) {
	return Maps(k.kv, prefix, offset)
}
