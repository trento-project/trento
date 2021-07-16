package runner

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	consultemplatelogging "github.com/hashicorp/consul-template/logging"
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

func NewTemplateRunner(runnerConfig *Config) (*manager.Runner, error) {
	consulConfig := consultemplateconfig.DefaultConsulConfig()
	consulConfig.Address = &runnerConfig.ConsulAddr

	loggingConfig := &consultemplatelogging.Config{
		Level:  strings.ToUpper(runnerConfig.ConsulTemplateLogLevel),
		Syslog: false,
		Writer: os.Stdout,
	}

	consultemplatelogging.Setup(loggingConfig)

	cTemplateConfig := consultemplateconfig.DefaultConfig()
	cTemplateConfig.Consul = consulConfig

	contents := ansibleHostsTemplate
	destination := path.Join(runnerConfig.AnsibleFolder, ansibleHostFile)
	*cTemplateConfig.Templates = append(
		*cTemplateConfig.Templates,
		&consultemplateconfig.TemplateConfig{
			Contents:    &contents,
			Destination: &destination,
		},
	)

	cTemplateConfig.Finalize()

	cTemplateRunner, err := manager.NewRunner(cTemplateConfig, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not start consul-template")
	}

	return cTemplateRunner, nil
}

func (c *Runner) startConsulTemplate(renderedWg *sync.WaitGroup) {
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

func (c *Runner) stopConsulTemplate() {
	log.Println("Stopping consul-template")
	c.templateRunner.StopImmediately()
	log.Println("Stopped consul-template")
}
