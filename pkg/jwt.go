package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

func GenerateJWT(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}
