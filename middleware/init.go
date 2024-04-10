package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RouterInit() cors.Config {
	// set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                        // allowed origin，use * represent for plural
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // allowed http method
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // allowed http header
	return config
}

func LimitCount(context *gin.Context) {
	limiter := rate.NewLimiter(200, 1)
	if !limiter.Allow() {
		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
		return
	}
	context.Next()
}
