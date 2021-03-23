module github.com/trento-project/trento

go 1.16

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/consul/api v1.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
)

replace github.com/trento-project/trento => ./
