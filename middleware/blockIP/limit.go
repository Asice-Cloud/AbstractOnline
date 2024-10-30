package blockIP

import (
	"Abstract/config"
	"Abstract/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// LimitCount check the visit frequency if it is too frequent, blocking the IP
func LimitCount(context *gin.Context) (err string) {
	ip := context.ClientIP()
	limiter := rate.NewLimiter(1000, 1)
	if !limiter.Allow() {
		// add this ip into blocked ip
		mu.Lock()
		err := AddBlockIP(ip)
		if err != nil {
			return ""
		}
		mu.Unlock()
		response.RespError(context, response.CodeServerBusy)
	}
	return ""
}

// BlockIPMiddleware the middleware to block malicious ip
func BlockIPMiddleware(context *gin.Context) {
	ip := context.ClientIP()
	checkResponse := LimitCount(context)
	if checkResponse != "" {
		response.RespErrorWithMsg(context, response.CodeInvalidAuthFormat, checkResponse)
		context.Abort()
		return
	}
	// Check if the IP is blocked
	val, err := config.Rdb.Get(context, ip).Result()
	if err == redis.Nil {
		// Key does not exist, continue to the next middleware
		context.Next()
		return
	} else if err != nil {
		// An actual error occurred, return a 503 error
		response.RespError(context, response.CodeServerBusy)
		context.Abort()
		return
	}
	if val == "blocked" {
		response.RespErrorWithMsg(context, response.CodeInvalidAuthFormat, errors.New("ip banned"))
		context.Abort()
		return
	}
	context.Next()
}
