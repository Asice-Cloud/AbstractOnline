package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySecret = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

const TokenExpireDuration = time.Hour * 2

var commonIssuer, commonSub = "CHATSYSTEM", "CHATONLINE"

type ClaimsFunc interface {
	jwt.Claims
	GetUserID() (int, error)
	GetUserName() (string error)
}

type Claims struct {
	jwt.Claims
	userName string
	userID   int
}

func (c Claims) GetUserID() (int, error) {
	return c.userID, nil
}

func (c Claims) GetUserName() (string, error) {
	return c.userName, nil
}

// 本来不想写，但是用多了还是要写
func KeyFunc(_ *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": commonIssuer,
		"aud": username,
		"sub": commonSub,
		"exp": time.Now().Add(TokenExpireDuration).Unix(),
	})
	aToken, _ = token.SignedString(mySecret)

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": commonIssuer,
		"sub": commonSub,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})
	rToken, err = rt.SignedString(mySecret)
	return
}

func ParseToken(tokenString string) (Claims, error) {
	var res Claims

	token, err := jwt.ParseWithClaims(tokenString, res, KeyFunc)
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return res, err
}

func RefreshToken(tokenString string) (newAToken, newRToken string, err error) {
	token, err := jwt.Parse(tokenString, KeyFunc)
	if !token.Valid {
		return "", "", errors.New("invalid token")
	}
	if err != nil {
		return "", "", errors.New("wrong when prase token")
	}
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return "", "", errors.New("invalid token")
	}
	nt := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": commonIssuer,
		"sub": commonSub,
		"aud": aud,
		"exp": time.Now().Add(TokenExpireDuration).Unix(),
	})
	newAToken, _ = nt.SignedString(mySecret)

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": commonIssuer,
		"sub": commonSub,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})
	newRToken, err = rt.SignedString(mySecret)
	return
}
