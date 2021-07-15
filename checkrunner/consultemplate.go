package checkrunner

import (
	"fmt"
	"path"
	"sync"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/internal/consul"
)

var ansibleHostsTemplate = fmt.Sprintf(`{{- with node }}
{{- $nodename := .Node.Node }}
[all]
{{- range nodes }}
{{- if ne .Node $nodename }}
{{ .Node }}
{{- end }}
{{- end }}
{{- end }}
{{- range $key, $pairs := tree "%[1]s" | byKey }}
[{{ key (print "%[1]s" $key "/name") }}]
{{- range tree (print "%[1]s" $key "/crmmon/Nodes") }}
{{- if .Key | contains "/Name" }}
{{ .Value }}
{{- end }}
{{- end }}
{{- end }}
`, consul.KvClustersPath)

const ansibleHostFile = "ansible_hosts"

func NewTemplateRunner(folder, consuldAddr string) (*manager.Runner, error) {
	config := consultemplateconfig.DefaultConfig()
	consulConfig := consultemplateconfig.DefaultConsulConfig()
	consulConfig.Address = &consuldAddr
	config.Consul = consulConfig

	contents := ansibleHostsTemplate
	destination := path.Join(folder, ansibleHostFile)
	*config.Templates = append(
		*config.Templates,
		&consultemplateconfig.TemplateConfig{
			Contents:    &contents,
			Destination: &destination,
		},
	)

	runner, err := manager.NewRunner(config, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not start consul-template")
	}

	return runner, nil
}

func (c *CheckRunner) startConsulTemplate(renderedWg *sync.WaitGroup) {
	var rendered bool = false
	go c.templateRunner.Start()
	defer c.stopConsulTemplate()

	for {
		select {
		case <-c.templateRunner.TemplateRenderedCh():
			if rendered {
				continue
			}
			log.Info("Template rendered and file created")
			renderedWg.Done()
			rendered = true
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *CheckRunner) stopConsulTemplate() {
	log.Println("Stopping consul-template")
	c.templateRunner.StopImmediately()
	log.Println("Stopped consul-template")
}
