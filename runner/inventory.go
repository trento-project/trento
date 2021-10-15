package runner

import (
  "os"
  "fmt"
  "text/template"

  "github.com/trento-project/trento/internal/consul"
  "github.com/trento-project/trento/internal/cluster"
  "github.com/trento-project/trento/internal/hosts"
  "github.com/trento-project/trento/internal/cloud"
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
{{ .Name }} ansible_host={{ .AnsibleHost }} ansible_user={{ .AnsibleUser }} {{ range $key, $value := .Variables }}{{ $key }}={{ $value }}{{ end }}
{{- end }}
{{- range .Groups }}
[{{ .Name }}]
{{- range .Nodes }}
{{ .Name }} ansible_host={{ .AnsibleHost }} ansible_user={{ .AnsibleUser }} {{ range $key, $value := .Variables }}{{ $key }}={{ $value }}{{ end }}
{{- end }}
{{- end }}
`
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

func NewClusterInventoryContent(client consul.Client) (*InventoryContent, error) {
  content := &InventoryContent{}

  clusters, err := cluster.Load(client)
  if err != nil {
    return nil, err
  }

  for clusterId, clusterData := range clusters {
    nodes := []*Node{}
    for _, node := range clusterData.Crmmon.Nodes {
      node := &Node{
        Name: node.Name,
        AnsibleUser: DefaultUser,
        Variables: make(map[string]interface{}),
      }

      // Get checks
      checks, err := cluster.GetCheckSelection(client, clusterId)
      if err == nil {
        node.Variables[clusterSelectedChecks] = checks
      }

      // Get ansible host
      hostList, err := hosts.Load(client, fmt.Sprintf("Node == \"%s\"", node.Name), []string{})
      if err == nil && len(hostList) > 0 {
        node.AnsibleHost = hostList[0].Node.Address
      }

      // Get ansible user
      cloudData, err := cloud.Load(client, node.Name)
      if err == nil {
        switch cloudData.Provider {
        case cloud.Azure:
          azureData := cloudData.Metadata.(cloud.AzureMetadata)
          node.AnsibleUser = azureData.Compute.OsProfile.AdminUserName
        }
      }

      nodes = append(nodes, node)
    }
    group := &Group{Name: clusterId, Nodes: nodes}

    content.Groups = append(content.Groups, group)
  }

  return content, nil
}
