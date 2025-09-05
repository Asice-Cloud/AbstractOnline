package router

import (
	"Abstract/controller"
	am "Abstract/controller/admin_module"
	az "Abstract/controller/authorization"
	um "Abstract/controller/user_module"
	vc "Abstract/controller/verification"
	"Abstract/middleware/debug"

	"github.com/gin-gonic/gin"
)

func Routers(router *gin.Engine) {
	// Apply debug defense middleware globally to all routes
	router.Use(debug.DebugDefenseMiddleware())

	router.GET("/index", controller.Welcome)

	// Debug detection API endpoint
	router.POST("/api/debug-detected", debug.DebugDetectionHandler())

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
