package runner

import (
	"fmt"
	"os"
	"text/template"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

type InventoryContent struct {
	Groups []*Group
	Nodes  []*Node
}

type Group struct {
	Name  string
	Nodes []*Node
}

type Node struct {
	Name        string
	AnsibleHost string
	AnsibleUser string
	Variables   map[string]interface{}
}

const (
	inventoryTemplate = `{{- range .Nodes }}
{{ .Name }} ansible_host={{ .AnsibleHost }} ansible_user={{ .AnsibleUser }} {{ range $key, $value := .Variables }}{{ $key }}={{ $value }} {{ end }}
{{- end }}
{{- range .Groups }}
[{{ .Name }}]
{{- range .Nodes }}
{{ .Name }} ansible_host={{ .AnsibleHost }} ansible_user={{ .AnsibleUser }} {{ range $key, $value := .Variables }}{{ $key }}={{ $value }} {{ end }}
{{- end }}
{{- end }}
`
	DefaultUser           string = "root"
	clusterSelectedChecks string = "cluster_selected_checks"
)

func CreateInventory(destination string, content *InventoryContent) error {
	t := template.Must(template.New("").Parse(inventoryTemplate))

	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	err = t.Execute(f, content)
	if err != nil {
		return nil
	}
	f.Close()

	return nil
}

// Local methods created to make the mocking possible
// These methods will be replaced once we have the new backend, so bear with them
var getClusters = func(client consul.Client) (map[string]*cluster.Cluster, error) {
	clusters, err := cluster.Load(client)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

var getCheckSelection = func(client consul.Client, clusterId string) (string, error) {
	checks, err := cluster.GetCheckSelection(client, clusterId)
	if err != nil {
		return "", err
	}

	return checks, nil
}

var getNodeAddress = func(client consul.Client, node string) (string, error) {
	hostList, err := hosts.Load(client, fmt.Sprintf("Node == \"%s\"", node), []string{})
	if err == nil && len(hostList) > 0 {
		return hostList[0].Node.Address, nil
	}

	return "", err
}

var getConnectionName = func(client consul.Client, clusterId string, node string) (string, error) {
	connectionData, err := cluster.GetConnectionSettings(client, clusterId)
	if err != nil {
		return "", err
	}

	user, found := connectionData[node]
	if !found {
		return "", fmt.Errorf("connection data for %s node not found", node)
	}

	return user.(string), nil
}

var getCloudUserName = func(client consul.Client, node string) (string, error) {
	cloudData, err := cloud.Load(client, node)
	if err != nil {
		return "", err
	}

	switch cloudData.Provider {
	case cloud.Azure:
		azureData := cloudData.Metadata.(cloud.AzureMetadata)
		return azureData.Compute.OsProfile.AdminUserName, nil
	default:
		return DefaultUser, nil
	}
}

func NewClusterInventoryContent(client consul.Client) (*InventoryContent, error) {
	content := &InventoryContent{}

	clusters, err := getClusters(client)
	if err != nil {
		return nil, err
	}

	for clusterId, clusterData := range clusters {
		nodes := []*Node{}
		for _, node := range clusterData.Crmmon.Nodes {
			node := &Node{
				Name:        node.Name,
				AnsibleUser: DefaultUser,
				Variables:   make(map[string]interface{}),
			}

			checks, err := getCheckSelection(client, clusterId)
			if err == nil {
				node.Variables[clusterSelectedChecks] = checks
			}

			address, err := getNodeAddress(client, node.Name)
			if err == nil {
				node.AnsibleHost = address
			}

			userName, err := getConnectionName(client, clusterId, node.Name)
			if err == nil {
				node.AnsibleUser = userName
			}

			// if the node user name is not provided by the user, get the cloud data
			if err != nil || len(userName) == 0 {
				cloudUser, err := getCloudUserName(client, node.Name)
				if err == nil {
					node.AnsibleUser = cloudUser
				}
			}

			nodes = append(nodes, node)
		}
		group := &Group{Name: clusterId, Nodes: nodes}

		content.Groups = append(content.Groups, group)
	}

	return content, nil
}
