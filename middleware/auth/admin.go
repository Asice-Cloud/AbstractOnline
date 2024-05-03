package auth

import (
	"Chat/controller"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"net/http"
	"strings"
)

var jwtKey = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

func AdminAuth(ctx *gin.Context) {
	authHeader := controller.SessionGet(ctx, "token").(string)
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not logged in",
		})
		return
	}

	tokenString := strings.Split(authHeader, "Bearer ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := claims["name"].(string)
		if user != "admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Not logged in as admin",
			})
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return
	}

	// If we reach this point, the user is logged in as admin
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logged in as admin",
	})
}
