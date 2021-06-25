package agent

import (
	"path"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const nodeMetadataTemplate = `{
  "node_meta": {
    {{- $first := true }}
    {{- with node }}
    {{- $nodename := .Node.Node }}
    {{- range nodes }}
    {{- if eq .Node $nodename }}
    {{- range ls (print "trento/v0/hosts/" $nodename "/metadata") }}
      {{- if $first }}{{ $first = false }}{{ else }},{{ end }}
      "trento-{{ .Key }}": "{{ .Value }}"
    {{- end }}
    {{- end }}
    {{- end }}
    {{- end }}
  }
}`

const metaDataFile = "trento-metadata.json"

func NewTemplateRunner(consulConfigDir string) (*manager.Runner, error) {
	config := consultemplateconfig.DefaultConfig()
	contents := nodeMetadataTemplate
	destination := path.Join(consulConfigDir, metaDataFile)
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

func (a *Agent) startConsulTemplate() {
	go a.templateRunner.Start()
	defer a.stopConsulTemplate()

	for {
		select {
		case <-a.templateRunner.TemplateRenderedCh():
			log.Printf("Template rendered. Reloading agent configuration...")
			err := a.consul.Agent().Reload()
			if err != nil {
				log.Printf("Error reloading agent meta-data information: %s", err)
			} else {
				log.Print("Agent meta-data correctly reloaded")
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func (a *Agent) stopConsulTemplate() {
	log.Println("Stopping consul-template")
	a.templateRunner.StopImmediately()
	log.Println("Stopped consul-template")
}
