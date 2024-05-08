package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

var jwtKey = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

func AdminAuthMiddleware(ctx *gin.Context) {
	// Get the session
	session := sessions.Default(ctx)
	authToken := session.Get("admin")
	/*log.Printf("AdminAuthMiddleware: session: %v\n", sessions.Default(ctx)) // Add log here
	log.Printf("AdminAuthMiddleware: session token: %v\n", authToken)       // Add log here*/

	// Check if the session exists
	if authToken == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not logged in",
		})
		ctx.Abort()
		return
	}

	// Check if the session is valid
	tokenString, ok := authToken.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		ctx.Abort()
		return

	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Printf("AdminAuthMiddleware: error parsing token: %v\n", err) // Add log here
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		ctx.Abort()
		return
	}

	// Check if the user is an admin
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role := claims["role"].(string)
		if role != "admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Not logged in as admin",
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

	// If we reach this point, the user is logged in as admin
	ctx.Next()
}
