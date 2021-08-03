package web

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const AlertsKey string = "alerts"

type Alert struct {
	Type  string
	Title string
	Text  string
}

func (a *Alert) GetIcon() string {
	switch a.Type {
	case "success":
		return "check_circle"
	case "warning":
		return "warning"
	case "danger":
		return "error"
	default:
		return "info"
	}
}

func InitAlerts() {
	gob.Register(Alert{})
}

func StoreAlert(c *gin.Context, a *Alert) {
	session := sessions.Default(c)
	session.AddFlash(a, AlertsKey)
	session.Save()
}

func GetAlerts(c *gin.Context) []*Alert {
	session := sessions.Default(c)
	f := session.Flashes(AlertsKey)
	var alerts []*Alert

	for _, alertI := range f {
		alert, _ := alertI.(Alert)
		alerts = append(alerts, &alert)
	}
	session.Save()

	return alerts
}
