package middleware

import (
	"Chat/config"
	"Chat/middleware/auth"
	"Chat/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
	"sync"
)

var (
	// BlockIP blocked IP:
	// 1. use redis to store the IP
	// 2. use the IP as the key, and the attribute as value
	// 3. if the visit frequency is too high, block the IP
	// 4. if the IP is blocked, return 429 status code
	// 5. if the IP is not blocked, continue to the next middleware
	BlockIP = make(map[int]string)
	mu      sync.Mutex
)

// LimitCount check the visit frequency if it is too frequent, blocking the IP
func LimitCount(context *gin.Context) (err string) {
	ip := context.ClientIP()
	limiter := rate.NewLimiter(200, 1)
	if !limiter.Allow() {
		// add this ip into blocked ip
		mu.Lock()
		BlockIP[len(BlockIP)] = ip
		err := auth.AddBlockIP(BlockIP)
		if err != nil {
			return ""
		}
		mu.Unlock()
		return response.CustomError{Code: -1, Msg: "Too many requests"}.Error()
	}
	return ""
}

// BlockIPMiddleware the middleware to block malicious ip
func BlockIPMiddleware(context *gin.Context) {
	ip := context.ClientIP()
	checkResponse := LimitCount(context)
	if checkResponse != "" {
		context.JSON(429, checkResponse)
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
		context.JSON(503, response.CustomError{-1, "Service Unavailable"}.Error())
		context.Abort()
		return
	}
	if val == "blocked" {
		context.JSON(403, response.CustomError{-1, "Forbidden"}.Error())
		context.Abort()
		return
	}
	context.Next()
}