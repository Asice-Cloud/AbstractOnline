package router

import (
	"Chat/controller"
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
		userRouter.GET("/login", controller.Login)
		userRouter.GET("/index", controller.Index)
		userRouter.GET("/adduser", controller.CreateUser)
		userRouter.DELETE("/deluser", controller.DeleteUser)
		userRouter.POST("/updateuser", controller.UpdateUser)
		userRouter.DELETE("/logout", controller.Logout)
		userRouter.POST("/searchuser", controller.SearchUser)
	}

	gitRouter := router.Group("/git")
	{
		gitRouter.GET("/login", controller.GitLogin)
		gitRouter.GET("/callback", controller.GitCallBack)
	}

	//admin module
	adminRouter := router.Group("/admin")

	adminRouter.GET("/login", controller.AdminLogin)

	//adminRouter.Use(auth.AdminJwtAuthMiddleware())
	{
		adminRouter.GET("/userlist", controller.GetUserList)
		adminRouter.GET("/retrievalblockip", controller.BlockIPRetrieval)
		adminRouter.DELETE("/deleteblockip", controller.BlockIPRemove)
	}
}
