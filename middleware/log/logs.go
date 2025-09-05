package log

import (
	"Abstract/config"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	dc = config.Colors["debug"]
	ic = config.Colors["info"]
	wc = config.Colors["warn"]
	ec = config.Colors["error"]
	rc = config.Colors["reset"]
)

func isDebugLoggingEnabled() bool {
	return gin.Mode() == gin.DebugMode
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Make a copy of the context for the goroutine
		cCp := *c
		start := time.Now()
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		cCp.Writer = w

		var requestBody []byte
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = make([]byte, len(bodyBytes))
				copy(requestBody, bodyBytes)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset request body
			}
		}
		c.Next()
		// Use a goroutine for logging
		go func() {
			status := cCp.Writer.Status()
			path := cCp.Request.URL.Path
			query := cCp.Request.URL.RawQuery
			method := cCp.Request.Method
			clientIP := cCp.ClientIP()
			userAgent := cCp.Request.UserAgent()
			cost := time.Since(start)
			responseHeaders := cCp.Writer.Header()
			responseBody := w.body.Bytes()

			if isDebugLoggingEnabled() {
				requestHeaders, _ := httputil.DumpRequest(cCp.Request, false)
				config.Lg.Debug("Debug info: \n",
					zap.String("path", path),
					zap.String(fmt.Sprintf("%s%s%s", dc, "method", rc), method),
					zap.Int("status", status),
					zap.Any("requestHeaders", string(requestHeaders)),
					zap.ByteString("requestBody", requestBody),
					zap.Any("responseHeaders", responseHeaders),
					zap.ByteString("responseBody", responseBody),
				)
			} else {
				// fuck you golang, why can not override function
				// looking at these fucking poor function
				// who never will stand
				if status >= 200 && status < 400 {
					config.Lg.Info(path,
						zap.String("method", cCp.Request.Method),
						zap.String("query", query),
						zap.String("ip", clientIP),
						zap.String("user-agent", userAgent),
						zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
						zap.Duration("cost", cost),
					)
				} else if status >= 400 && status < 500 {
					config.Lg.Warn(path,
						zap.String("method", cCp.Request.Method),
						zap.String("query", query),
						zap.String("ip", cCp.ClientIP()),
						zap.String("user-agent", cCp.Request.UserAgent()),
						zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
						zap.Duration("cost", cost),
					)
				} else {
					config.Lg.Error(path,
						zap.String("method", cCp.Request.Method),
						zap.String("query", query),
						zap.String("ip", cCp.ClientIP()),
						zap.String("user-agent", cCp.Request.UserAgent()),
						zap.String("errors", cCp.Errors.ByType(gin.ErrorTypePrivate).String()),
						zap.Duration("cost", cost),
					)
				}
			}
		}()
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
					config.Lg.Fatal(c.Request.URL.Path,
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
