package auth

import (
	"Abstract/pkg"
	"Abstract/response"
	"Abstract/session"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

func AdminJwtAuthMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		authSession := session.SessionGet("user", c, "admin")
		if authSession == nil {
			response.RespErrorWithMsg(c, response.CodeInvalidToken, errors.New("empty token"))
			c.Abort()
			return
		}
		authToken, ok := authSession.(session.AdminSession)
		if !ok {
			response.RespError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}
		mc, err := pkg.ParseToken(authToken.AccessToken)
		if err != nil {
			log.Println(err)
			response.RespError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}
		if mc.UserID != 0 && mc.Username != "admin" {
			response.RespError(c, response.CodeInvalidToken)
			c.Abort()
			return
		} else {
			c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
		}
	}
}

func UserJwtAuthMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		authSession := session.SessionGet("user", c, "user")
		if authSession == nil {
			response.RespErrorWithMsg(c, response.CodeInvalidToken, errors.New("empty token"))
			c.Abort()
			return
		}
		authToken, ok := authSession.(session.UserSession)
		if !ok {
			response.RespError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}
		_, err := pkg.ParseToken(authToken.AccessToken)
		if err != nil {
			log.Println(err)
			response.RespError(c, response.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}
