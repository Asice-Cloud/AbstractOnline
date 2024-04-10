package Middleware

import (
	"fmt"
	"net/http"
	"time"

	"Chat/Config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var logPool = Config.NewLogPool(10)

// Limit the max count of synchronic requesets
func LimitCount(context *gin.Context) {
	limiter := rate.NewLimiter(200, 1)
	if !limiter.Allow() {
		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
		return
	}
	context.Next()
}

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
