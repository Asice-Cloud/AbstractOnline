package auth

import (
	"Chat/controller"
	"Chat/pkg"
	"Chat/response"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

func AdminJwtAuthMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		authSession := controller.SessionGet(c, "admin")
		if authSession == nil {
			response.RespErrorWithMsg(c, response.CodeInvalidToken, errors.New("empty token"))
			c.Abort()
			return
		}
		authToken, ok := authSession.(controller.AdminSession)
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
		authSession := controller.SessionGet(c, "user")
		if authSession == nil {
			response.RespErrorWithMsg(c, response.CodeInvalidToken, errors.New("empty token"))
			c.Abort()
			return
		}
		authToken, ok := authSession.(controller.UserSession)
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
