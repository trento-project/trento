package web

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const AlertsKey string = "alerts"

// Here some generic alerts
var AlertCatalogNotFound = func() Alert {
	return Alert{
		Type:  "danger",
		Title: "Error loading the checks catalog",
		Text: "Checks catalog couldn't be retrieved. Check if the ARA service is running" +
			" and the --ara-addr flag is pointing corretly to the service",
	}
}

var AlertConnectionDataNotFound = func() Alert {
	return Alert{
		Type:  "danger",
		Title: "Error loading the connection data",
		Text: "Connection data couldn't be retrieved.",
	}
}

var CheckResultsNotFound = func() Alert {
	return Alert{
		Type:  "danger",
		Title: "Error loading the checks result",
		Text:  "Checks result couldn't be retrieved. Check if the Trento runner is running",
	}
}

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

func StoreAlert(c *gin.Context, a Alert) {
	session := sessions.Default(c)
	session.AddFlash(a, AlertsKey)
	session.Save()
}

func GetAlerts(c *gin.Context) []Alert {
	session := sessions.Default(c)
	f := session.Flashes(AlertsKey)
	var alerts []Alert

	for _, alertI := range f {
		alert, _ := alertI.(Alert)
		alerts = append(alerts, alert)
	}
	session.Save()

	return alerts
}
