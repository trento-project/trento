module github.com/trento-project/trento

go 1.16

require (
	github.com/ClusterLabs/ha_cluster_exporter v0.0.0-20210420075709-eb4566acab09
	github.com/SUSE/sap_host_exporter v0.0.0-20210426144122-68bbf2f1b490
	github.com/aquasecurity/bench-common v0.4.4
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/mock v1.4.3
	github.com/gomarkdown/markdown v0.0.0-20210514010506-3b9f47219fe7
	github.com/hashicorp/consul-template v0.25.2
	github.com/hashicorp/consul/api v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.1.2
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tdewolff/minify/v2 v2.9.16
)

replace github.com/trento-project/trento => ./
