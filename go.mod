module github.com/trento-project/trento

go 1.16

require (
	github.com/aquasecurity/bench-common v0.4.4
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/consul-template v0.25.2
	github.com/hashicorp/consul/api v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tdewolff/minify/v2 v2.9.16
)

replace github.com/trento-project/trento => ./
