package router

import (
	"Abstract/controller"
	am "Abstract/controller/admin_module"
	az "Abstract/controller/authorization"
	um "Abstract/controller/user_module"
	vc "Abstract/controller/verification"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(router *gin.Engine) {
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/index", controller.Welcome)

	userRouter := router.Group("/user")
	//userRouter.Use(blockIP.BlockIPMiddleware)
	//userRouter.Use(auth.UserJwtAuthMiddleware())
	{
		userRouter.GET("/login", um.Login)
		userRouter.GET("/index", um.Index)
		userRouter.GET("/adduser", um.CreateUser)
		userRouter.DELETE("/deluser", um.DeleteUser)
		userRouter.POST("/updateuser", um.UpdateUser)
		userRouter.DELETE("/logout", um.Logout)
		userRouter.POST("/searchuser", um.SearchUser)

		userRouter.GET("/before", um.Before)
		userRouter.GET("/home", um.Home)
		userRouter.GET("/ws", um.Ws)
	}

	av := router.Group("/av")
	{
		//Github Oauth
		av.GET("/login", az.GitLogin)
		av.GET("/callback", az.GitCallBack)

		// 验证滑块验证码
		av.GET("background", vc.GetBackground)
		av.GET("slider", vc.Slider)
		av.POST("verify", vc.Verify)
	}

	//admin module
	adminRouter := router.Group("/admin")

	adminRouter.GET("/login", am.AdminLogin)

	//adminRouter.Use(auth.AdminJwtAuthMiddleware())
	{
		adminRouter.GET("/userlist", am.GetUserList)
		adminRouter.GET("/retrievalblockip", am.BlockIPRetrieval)
		adminRouter.DELETE("/deleteblockip", am.BlockIPRemove)
	}
}
