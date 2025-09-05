package debug

import (
	"Abstract/server"
	"bytes"
	"strings"

	"github.com/gin-gonic/gin"
)

// responseWriter wraps gin.ResponseWriter to capture HTML responses
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return len(b), nil
}

// DebugDefenseMiddleware injects debug defense script into all HTML responses
func DebugDefenseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get defense configuration
		config := server.DefaultDefenseConfig()

		// Skip if defense is disabled
		if !config.Enabled {
			c.Next()
			return
		}

		// Create a custom response writer to capture the response
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		// Check if the response is HTML
		contentType := c.Writer.Header().Get("Content-Type")
		if strings.Contains(contentType, "text/html") || contentType == "" {
			// Get the response body
			responseBody := blw.body.String()

			// Only inject if this is an HTML document
			if strings.Contains(responseBody, "<html") || strings.Contains(responseBody, "<!DOCTYPE") {
				// Generate debug defense script
				debugScript := server.GenerateDebugDetectionScript(config.RedirectURL)

				// Inject the script before closing </head> tag or at the beginning of <body>
				if strings.Contains(responseBody, "</head>") {
					responseBody = strings.Replace(responseBody, "</head>", debugScript+"\n</head>", 1)
				} else if strings.Contains(responseBody, "<body>") {
					responseBody = strings.Replace(responseBody, "<body>", "<body>\n"+debugScript, 1)
				} else if strings.Contains(responseBody, "<html>") {
					// If no head or body tags, inject after <html>
					responseBody = strings.Replace(responseBody, "<html>", "<html>\n"+debugScript, 1)
				}
			}

			// Write the modified response
			c.Writer.Header().Set("Content-Length", string(rune(len(responseBody))))
			_, _ = c.Writer.WriteString(responseBody)
		} else {
			// For non-HTML responses, write the original body
			_, _ = c.Writer.Write(blw.body.Bytes())
		}
	}
}

// DebugDetectionHandler handles debug detection API calls
func DebugDetectionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := server.DefaultDefenseConfig()

		// Log the detection attempt (variables are used for potential future logging)
		_ = c.ClientIP()
		_ = c.GetHeader("User-Agent")

		// You can add logging here if needed:
		// logger.Info("Debug tools detected", map[string]interface{}{
		//     "ip": clientIP,
		//     "user_agent": userAgent,
		//     "timestamp": time.Now(),
		// })

		// Handle the detection
		server.HandleDebugDetection(c.Writer, c.Request, config)
	}
}
