package auth

import "github.com/gin-gonic/gin"

func AdminAuth(ctx *gin.Context) {
	// check the token
	token := ctx.GetHeader("Authorization")
	if token != "admin" {
		ctx.JSON(403, gin.H{
			"message": "Forbidden",
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}
