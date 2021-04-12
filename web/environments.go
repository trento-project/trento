package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aquasecurity/bench-common/check"
	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

type Environment struct {
	Name  string
	Nodes []*Node
}

type EnvironmentList map[string]*Environment

type Node struct {
	consulApi.Node
	client consul.Client
}

func (n *Node) Health() string {
	checks, _, _ := n.client.Health().Node(n.Name(), nil)
	return checks.AggregatedStatus()
}

func (n *Node) Name() string {
	return n.Node.Node
}

// todo: this method was rushed, needs to be completely rewritten to have the checker webservice decoupled in a dedicated HTTP client
func (n *Node) Checks() *check.Controls {
	checks := &check.Controls{}

	var err error
	resp, err := http.Get(fmt.Sprintf("http://%s:%d", n.Address, 8700)) // todo: use a Consul service instead of directly addressing the node IP and port
	if err != nil {
		log.Print(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	err = json.Unmarshal(body, checks)
	if err != nil {
		log.Print(err)
		return nil
	}
	return checks
}

func NewEnvironmentsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
		if err != nil {
			_ = c.Error(err)
			return
		}
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
		environments[node.Datacenter].Nodes = append(environments[node.Datacenter].Nodes, &Node{*node, client})
	}

	return environments, nil
}
