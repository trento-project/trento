package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/web/services"
)

//go:generate swag init -g api.go -o ../docs/api
// @title Trento API
// @version 1.0
// @description Trento API

// @contact.name Trento Project
// @contact.url https://www.trento-project.io
// @contact.email  trento-project@suse.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
// @schemes http

func ApiPingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}


type JSONGroups []*JSONGroup

type JSONGroup struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func ApiGetPromHttpSd(s services.HostsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var groups JSONGroups

		hosts, err := s.GetAll(nil, nil)
		if err != nil {
			c.Error(err)
			return
		}

		for _, host := range hosts {
			group := &JSONGroup{
				Targets: []string{fmt.Sprintf("%s:9999", host.SSHAddress)},
				Labels: map[string]string{
					"hostname": host.Name,
				},
			}
			groups = append(groups, group)
		}

		c.JSON(http.StatusOK, groups)
	}
}
