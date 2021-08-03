module github.com/trento-project/trento

go 1.16

require (
	github.com/ClusterLabs/ha_cluster_exporter v0.0.0-20210420075709-eb4566acab09
	github.com/aquasecurity/bench-common v0.4.4
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/gomarkdown/markdown v0.0.0-20210514010506-3b9f47219fe7
	github.com/hashicorp/consul-template v0.25.2
	github.com/hashicorp/consul/api v1.4.0
	github.com/hooklift/gowsdl v0.5.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.1.2
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tdewolff/minify/v2 v2.9.16
	github.com/vektra/mockery/v2 v2.9.0
)

replace github.com/trento-project/trento => ./
