package safe

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeInputMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() == "application/x-www-form-urlencoded" {
			c.Request.ParseForm()
			sanitizer := bluemonday.StrictPolicy()
			for key, vals := range c.Request.PostForm {
				for i, val := range vals {
					c.Request.PostForm[key][i] = sanitizer.Sanitize(val)
				}
			}
		}
		c.Next()
	}
}
