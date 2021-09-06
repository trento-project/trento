package runner

import (
	"fmt"
	"os"
	"path"
	"strings"

	consultemplateconfig "github.com/hashicorp/consul-template/config"
	consultemplatelogging "github.com/hashicorp/consul-template/logging"
	"github.com/hashicorp/consul-template/manager"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/consul"
)

const clusterSelectedChecks string = "cluster_selected_checks"

var ansibleHostsTemplate = fmt.Sprintf(`
{{- range $clusterId, $pairs := tree "%[2]s" | byKey }}
[{{ key (print (printf "%[3]s" $clusterId) "/id") }}]
{{- range tree (print (printf "%[3]s" $clusterId) "/crmmon/Nodes") }}
{{- if .Key | contains "/Name" }}
{{- $nodename := .Value }}
{{- $user := keyOrDefault (print (printf "%[5]s" $clusterId) "/" $nodename) "root" }}
{{ $nodename }} %[1]s={{ keyOrDefault  (printf "%[4]s" $clusterId) "" }} ansible_host={{ range nodes }}{{ if eq .Node $nodename }}{{ .Address }}{{ end }}{{ end }} ansible_user={{ $user }}
{{- end }}
{{- end }}
{{- end }}
`,
	clusterSelectedChecks,
	consul.KvClustersPath,
	consul.KvClustersDiscoveredPath,
	consul.KvClustersChecksPath,
	consul.KvClustersConnectionPath)

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

	cTemplateConfig.Once = true

	cTemplateConfig.Finalize()

	cTemplateRunner, err := manager.NewRunner(cTemplateConfig, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not start consul-template")
	}

	return cTemplateRunner, nil
}

func (c *Runner) startConsulTemplate() {
	go c.templateRunner.Start()
	defer c.stopConsulTemplate()

	for {
		select {
		case <-c.templateRunner.TemplateRenderedCh():
			log.Info("Template rendered and file created")
			return
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
