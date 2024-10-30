package utils

import (
	"fmt"
	"time"
)
import "github.com/golang-jwt/jwt/v5"

var secret_key = []byte("abstract")

func CreateToken(un string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": un,
		"expire":   time.Now().Add(time.Hour * 24).Unix(),
	})
	token_string, err := token.SignedString(secret_key)
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func VerifyToken(token_string string) error {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})
	if err != nil {
		return err
	}
	if token.Valid {
		return fmt.Errorf("token is invalid")
	}
	return nil
}
