package web

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const AlertsKey string = "alerts"

type Alert struct {
	Type  string
	Title string
	Text  string
}

func (a *Alert) String() string {
	return fmt.Sprintf("%s::%s::%s", a.Type, a.Title, a.Text)
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

func alertStrToStruct(str string) (*Alert, error) {
	items := strings.Split(str, "::")
	if len(items) != 3 {
		return nil, fmt.Errorf("Malformed string. The string must have the type, title and text elements splitted by ::")
	}

	return &Alert{
		Type:  items[0],
		Title: items[1],
		Text:  items[2],
	}, nil
}

func StoreAlert(c *gin.Context, a *Alert) {
	session := sessions.Default(c)
	session.AddFlash(a.String(), AlertsKey)
	session.Save()
}

func GetAlerts(c *gin.Context) []*Alert {
	session := sessions.Default(c)
	f := session.Flashes(AlertsKey)
	var alerts []*Alert

	for _, alertI := range f {
		alert, _ := alertStrToStruct(alertI.(string))
		alerts = append(alerts, alert)
	}
	session.Save()

	return alerts
}
