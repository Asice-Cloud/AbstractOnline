package router

import (
	"Chat/controller"
	"Chat/controller/admin_module"
	"Chat/controller/authorization"
	"Chat/controller/user_module"
	"Chat/controller/verification"
	"Chat/middleware/blockIP"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(router *gin.Engine) {
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/index", controller.Welcome)

	userRouter := router.Group("/user")
	userRouter.Use(blockIP.BlockIPMiddleware)
	//userRouter.Use(auth.UserJwtAuthMiddleware())
	{
		userRouter.GET("/login", user_module.Login)
		userRouter.GET("/index", user_module.Index)
		userRouter.GET("/adduser", user_module.CreateUser)
		userRouter.DELETE("/deluser", user_module.DeleteUser)
		userRouter.POST("/updateuser", user_module.UpdateUser)
		userRouter.DELETE("/logout", user_module.Logout)
		userRouter.POST("/searchuser", user_module.SearchUser)

		userRouter.GET("/before", user_module.Before)
		userRouter.GET("/home", user_module.Home)
		userRouter.GET("/ws", user_module.Ws)
	}

	av := router.Group("/av")
	{
		//Github Oauth
		av.GET("/login", authorization.GitLogin)
		av.GET("/callback", authorization.GitCallBack)

		// 验证滑块验证码
		av.GET("background", verification.GetBackground)
		av.GET("slider", verification.Slider)
		av.POST("verify", verification.Verify)
	}

	//admin module
	adminRouter := router.Group("/admin")

	adminRouter.GET("/login", admin_module.AdminLogin)

	//adminRouter.Use(auth.AdminJwtAuthMiddleware())
	{
		adminRouter.GET("/userlist", admin_module.GetUserList)
		adminRouter.GET("/retrievalblockip", admin_module.BlockIPRetrieval)
		adminRouter.DELETE("/deleteblockip", admin_module.BlockIPRemove)
	}
}
