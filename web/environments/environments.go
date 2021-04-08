package envronments

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Environment struct {
	Name string
}

type EnvironmentList []Environment

func ListHandler(c *gin.Context) {
	var environments = EnvironmentList{
		Environment{
			Name: "test1",
		},
		Environment{
			Name: "test2",
		},
	}

	c.HTML(http.StatusOK, "environments.html.tmpl", gin.H{"Environments": environments})
}
