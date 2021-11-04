package web

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewLogHandler(instance string, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		var level log.Level

		switch true {
		case c.Writer.Status() >= 400 && c.Writer.Status() < 500:
			level = log.WarnLevel
		case c.Writer.Status() >= 500:
			level = log.ErrorLevel
		default:
			level = log.InfoLevel
		}

		logger.WithFields(log.Fields{
			"instance": instance,
			"client":   c.ClientIP(),
			"method":   c.Request.Method,
			"status":   c.Writer.Status(),
			"uri":      c.Request.RequestURI,
			"latency":  endTime.Sub(startTime),
		}).Log(level, "HTTP request")
	}
}
