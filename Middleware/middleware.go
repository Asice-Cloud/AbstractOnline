package Middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func LimitCount(context *gin.Context) {
	var limiter = rate.NewLimiter(200, 1)
	if !limiter.Allow() {
		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
		return
	}
	context.Next()
}
