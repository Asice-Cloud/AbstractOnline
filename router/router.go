package router

import (
	"Abstract/controller"
	am "Abstract/controller/admin_module"
	az "Abstract/controller/authorization"
	um "Abstract/controller/user_module"
	vc "Abstract/controller/verification"

	"github.com/gin-gonic/gin"
)

func Routers(router *gin.Engine) {
	router.GET("/index", controller.Welcome)

	userRouter := router.Group("/user")
	//userRouter.Use(blockIP.BlockIPMiddleware)
	{
		userRouter.GET("/index", um.Index)
		userRouter.GET("/before", um.Before)
		userRouter.GET("/home", um.Home)
		userRouter.GET("/ws", um.Ws)
	}

	av := router.Group("/av")
	{
		//GitHub Oauth
		av.GET("/login", az.GitLogin)
		av.GET("/callback", az.GitCallBack)

		// 验证滑块验证码
		av.GET("background", vc.GetBackground)
		av.GET("slider", vc.Slider)
		av.POST("verify", vc.Verify)
	}

	//admin module
	adminRouter := router.Group("/admin")

	{
		adminRouter.GET("/retrievalblockip", am.BlockIPRetrieval)
		adminRouter.DELETE("/deleteblockip", am.BlockIPRemove)
	}
}
