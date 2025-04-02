package auth

import (
	"Abstract/utils"
	"github.com/gin-gonic/gin"
	"testing"
)

type usertest struct {
	username string
	password string
}

func TestJwt(t *testing.T) {
	// Test jwt middleware
	r := gin.Default()
	r.Use(JWTmiddleware)
	r.POST("/test_get_token", func(c *gin.Context) {
		var user usertest
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid request",
			})
			return
		}
		token_string, err := utils.CreateToken(user.username)
		c.JSON(200, gin.H{
			"message": "success",
			"token":   token_string,
		})
	})
	r.GET("/test_verify", func(c *gin.Context) {
		token_string := c.GetHeader("Authorization")
		if token_string == "" {
			c.JSON(401, gin.H{
				"message": "Missing Authorization token",
			})
			return
		}
		err := utils.VerifyToken(token_string)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Invalid token",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
}
