package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySecret = []byte("HelloWorldCongratulationsYouHaveFoundTheSecretKey")

const TokenExpireDuration = time.Hour * 2

var commonIssuer, commonSub = "CHATSYSTEM", "CHATONLINE"

// 把人逼急了只会手动实现一遍Claims
type Claims struct {
	Issuer     string
	IssuedAt   time.Time
	Audience   []string
	Subject    string
	ExpireTime time.Time
	Username   string
	UserID     uint64
	NotBefore  time.Time
}

func (m Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{m.ExpireTime}, nil
}

func (m Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{m.IssuedAt}, nil
}

func (m Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{m.NotBefore}, nil
}

func (m Claims) GetIssuer() (string, error) {
	return m.Issuer, nil
}

func (m Claims) GetSubject() (string, error) {
	return m.Subject, nil
}

func (m Claims) GetAudience() (jwt.ClaimStrings, error) {
	return m.Audience, nil
}

func (c Claims) GetUserID() (uint64, error) {
	return c.UserID, nil
}

func (c Claims) GetUserName() (string, error) {
	return c.Username, nil
}

// 本来不想写，但是用多了还是要写
func KeyFunc(_ *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

func GenToken(userID uint64, username string) (aToken, rToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, Claims{
		Issuer:     commonIssuer,
		IssuedAt:   time.Now(),
		Subject:    commonSub,
		ExpireTime: time.Now().Add(TokenExpireDuration),
		Username:   username,
		UserID:     userID,
		NotBefore:  time.Now(),
	})
	aToken, _ = token.SignedString(mySecret)

	rt := jwt.NewWithClaims(jwt.SigningMethodES256, Claims{
		Issuer:     commonIssuer,
		Subject:    commonSub,
		ExpireTime: time.Now().Add(30 * time.Minute),
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
