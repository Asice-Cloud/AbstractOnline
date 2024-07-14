package log

import (
	"Chat/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Make a copy of the context for the goroutine
		cCp := c.Copy()
		start := time.Now()

		// Use a goroutine for logging
		go func() {
			path := cCp.Request.URL.Path
			query := cCp.Request.URL.RawQuery
			status := cCp.Writer.Status()
			cost := time.Since(start)

			if status >= 200 && status < 400 {
				zap.L().Info(path,
					zap.String("method", cCp.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", cCp.ClientIP()),
					zap.String("user-agent", cCp.Request.UserAgent()),
					zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
					zap.Duration("cost", cost),
				)
			} else if status >= 400 && status < 500 {
				zap.L().Warn(path,
					zap.String("method", cCp.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", cCp.ClientIP()),
					zap.String("user-agent", cCp.Request.UserAgent()),
					zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
					zap.Duration("cost", cost),
				)
			} else {
				zap.L().Error(path,
					zap.String("method", cCp.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", cCp.ClientIP()),
					zap.String("user-agent", cCp.Request.UserAgent()),
					zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
					zap.Duration("cost", cost),
				)
			}
		}()

		// Continue with the next middleware/handler
		c.Next()
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					config.Lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
