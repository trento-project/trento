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

const (
	DefaultUser           string = "root"
	clusterSelectedChecks string = "cluster_selected_checks"
)

var ansibleHostsTemplate = fmt.Sprintf(`
{{- /* Loop through the discovered clusters */}}
{{- range $clusterId, $pairs := tree "%[2]s" | byKey }}
[{{ key (print (printf "%[3]s" $clusterId) "/id") }}]
{{- range tree (print (printf "%[3]s" $clusterId) "/crmmon/Nodes") }}
{{- if .Key | contains "/Name" }}
{{- $nodename := .Value }}
{{- /* Get the node host address */}}
{{- $host := "" }}
{{- range nodes }}{{ if eq .Node $nodename }}{{ $host = .Address }}{{ end }}{{ end }}
{{- /* Get SSH connection username */}}
{{- $user := keyOrDefault (print (printf "%[5]s" $clusterId) "/" $nodename) "" }}
{{- /* If the user is not set, fallback to default values */}}
{{- if eq $user "" }}
{{- $cloudata := printf "%[6]s" $nodename }}
{{- $provider := keyOrDefault (print $cloudata "provider") "" }}
{{- if eq $provider "azure" }}
{{- $user = keyOrDefault (print $cloudata "metadata/compute/osprofile/adminusername") (printf "%[7]s") }}
{{- else }}
{{- $user = printf "%[7]s" }}
{{- end }}
{{- end }}
{{- /* Render the node entry */}}
{{ $nodename }} %[1]s={{ keyOrDefault  (printf "%[4]s" $clusterId) "" }} ansible_host={{ $host }} ansible_user={{ $user }}
{{- end }}
{{- end }}
{{- end }}
`,
	clusterSelectedChecks,
	consul.KvClustersPath,
	consul.KvClustersDiscoveredPath,
	consul.KvClustersChecksPath,
	consul.KvClustersConnectionPath,
	consul.KvHostsClouddataPath,
	DefaultUser,
)

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
