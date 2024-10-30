package auth

import (
	"Abstract/utils"
	"github.com/gin-gonic/gin"
)

func JWTmiddleware(c *gin.Context) {
	token_string := c.GetHeader("Authorization")
	if token_string == "" {
		c.JSON(401, gin.H{
			"message": "Missing Authorization token",
		})
		return
	}

	token_string = token_string[len("Bearer "):]

	err := utils.VerifyToken(token_string)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid token",
		})
		return
	}

	c.Next()
}
