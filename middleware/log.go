package middleware

import (
	"fmt"
	"time"

	"Chat/config"

	"github.com/gin-gonic/gin"
)

var logPool = config.NewLogPool(10)

// Limit the max count of synchronic requesets

// record the log of request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before request
		userIP := c.ClientIP()
		logTime := time.Now()
		logMessage := fmt.Sprintf("User IP: %s, Request Time: %s", userIP, logTime)
		logPool.Log(logMessage)

		c.Next()

		// After request
		status := c.Writer.Status()
		logMessage = fmt.Sprintf("User IP: %s, Request Time: %s, Status: %d", userIP, logTime, status)
		logPool.Log(logMessage)
	}
}
