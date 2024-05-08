package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

var mySecret = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

const TokenExpireDuration = time.Hour * 2

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func KeyFunc(_ *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	// own claim
	c := Claims{
		userID,
		"username",
		jwt.StandardClaims{ // 7 official sequence in JWT
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // expire time
			Issuer: "CHATSYSTEM", // register
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(), // 过期时间
		Issuer:    "CHATSYSYEM",                            // 签发人
	}).SignedString(mySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}

func ParseToken(tokenString string) (claims *Claims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(Claims)
	token, err = jwt.ParseWithClaims(tokenString, claims, KeyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, KeyFunc); err != nil {
		return
	}

	// 从旧access token中解析出claims数据	解析出payload负载信息
	var claims Claims
	_, err = jwt.ParseWithClaims(aToken, &claims, KeyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
