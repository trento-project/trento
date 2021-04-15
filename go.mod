module github.com/trento-project/trento

go 1.16

require (
	github.com/aquasecurity/bench-common v0.4.4
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/consul/api v1.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tdewolff/minify/v2 v2.9.16
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
)

replace github.com/trento-project/trento => ./
