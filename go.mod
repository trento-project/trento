module github.com/trento-project/trento

go 1.16

require (
	github.com/avast/retry-go/v4 v4.0.3
	github.com/gin-contrib/sessions v0.0.4
	github.com/gin-gonic/gin v1.7.7
	github.com/gomarkdown/markdown v0.0.0-20210514010506-3b9f47219fe7
	github.com/google/uuid v1.3.0
	github.com/hooklift/gowsdl v0.5.0
	github.com/lib/pq v1.10.4
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.32.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.4.1
	github.com/swaggo/swag v1.8.0
	github.com/tdewolff/minify/v2 v2.10.0
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/ugorji/go v1.1.13 // indirect
	github.com/vektra/mockery/v2 v2.10.0
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gorm.io/datatypes v1.0.2
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.15
)

replace github.com/trento-project/trento => ./
