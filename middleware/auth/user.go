package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func UserAuthMiddleware(ctx *gin.Context) {
	// Get the session
	session := sessions.Default(ctx)
	authToken := session.Get("user")

	// Check if the session exists
	if authToken == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not logged in",
		})
		ctx.Abort()
		return
	}

	// Check if the session is valid
	tokenString := authToken.(string)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		ctx.Abort()
		return
	}

	// Check if the user is a user
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role := claims["role"].(string)
		if role != "user" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Not logged in as user",
			})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		ctx.Abort()
		return
	}

	// If we reach this point, the user is logged in as user
	ctx.Next()
}
