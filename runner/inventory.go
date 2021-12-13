package runner

import (
	"encoding/json"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/api"
)

type InventoryContent struct {
	Groups []*Group
	Nodes  []*Node // this seems unused
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

func NewClusterInventoryContent(trentoApi api.TrentoApiService) (*InventoryContent, error) {
	content := &InventoryContent{}

	clustersSettings, err := trentoApi.GetClustersSettings()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clustersSettings {
		nodes := []*Node{}

		jsonSelectedChecks, err := json.Marshal(cluster.SelectedChecks)
		if err != nil {
			log.Errorf("error marshalling the cluster %s selected checks: %s", cluster.ID, err)
			continue
		}

		for _, host := range cluster.Hosts {
			node := &Node{
				Name:        host.Name,
				AnsibleHost: host.Address,
				AnsibleUser: host.User,
				Variables:   make(map[string]interface{}),
			}

			node.Variables[clusterSelectedChecks] = string(jsonSelectedChecks)

			nodes = append(nodes, node)
		}
		group := &Group{Name: cluster.ID, Nodes: nodes}

		content.Groups = append(content.Groups, group)
	}

	return content, nil
}
