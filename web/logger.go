package web

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewLogHandler(instance string, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Since(start)

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
			"instance":   instance,
			"client":     c.ClientIP(),
			"method":     c.Request.Method,
			"status":     c.Writer.Status(),
			"uri":        c.Request.RequestURI,
			"latency":    stop,
			"user_agent": c.Request.UserAgent(),
		}).Log(level, "HTTP request")
	}
}
