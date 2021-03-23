package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/SUSE/console-for-sap-applications/internal/consul"
)

type Environment struct {
	Name  string
	Nodes []*consulApi.Node
}

type EnvironmentList map[string]*Environment

func NewEnvironmentsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
		if err != nil {
			c.Error(err)
			return
		}

		//client.Health().Node()

		c.HTML(http.StatusOK, "environments.html.tmpl", gin.H{"Environments": environments})
	}
}

func loadEnvironments(client consul.Client) (EnvironmentList, error) {
	var environments = EnvironmentList{}

	dcs, err := client.Catalog().Datacenters()
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for datacenters")
	}
	for _, dc := range dcs {
		environments[dc] = &Environment{
			Name: dc,
		}
	}

	nodes, _, err := client.Catalog().Nodes(nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for nodes")
	}
	for _, node := range nodes {
		environments[node.Datacenter].Nodes = append(environments[node.Datacenter].Nodes, node)
	}

	return environments, nil
}
