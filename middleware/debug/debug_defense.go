package debug

import (
	"Abstract/server"
	"bytes"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// responseWriter wraps gin.ResponseWriter to capture HTML responses
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	// Capture written bytes into buffer for inspection/modification later.
	_, _ = w.body.Write(b)
	// Do not write through here; middleware will write the final bytes once it
	// decides whether to inject the script or pass the response unchanged.
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

		// Check if the response is HTML (only when Content-Type explicitly indicates HTML)
		// Note: handlers set headers on the wrapper (blw), so read from blw.Header()
		contentType := blw.Header().Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
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

			// Copy headers from the wrapper to the underlying writer so Content-Type and
			// other headers set by handlers are preserved.
			for k, vals := range blw.Header() {
				for _, v := range vals {
					blw.ResponseWriter.Header().Add(k, v)
				}
			}
			// Overwrite Content-Length with the modified body length
			blw.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
			_, _ = blw.ResponseWriter.Write([]byte(responseBody))
		} else {
			// For non-HTML responses, forward headers and the original buffered bytes
			// to the underlying ResponseWriter so headers (Content-Type, etc.) set
			// by the handler are preserved.
			for k, vals := range blw.Header() {
				for _, v := range vals {
					blw.ResponseWriter.Header().Add(k, v)
				}
			}
			_, _ = blw.ResponseWriter.Write(blw.body.Bytes())
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
