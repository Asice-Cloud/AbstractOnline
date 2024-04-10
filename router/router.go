package router

import (
	"Chat/controller"
	"Chat/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(router *gin.Engine) {
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/index", Service.Welcome)

	userRouter := router.Group("/user")
	userRouter.Use(middleware.BlockIPMiddleware)
	{
		userRouter.GET("/login", Service.Login)
		userRouter.GET("/index", Service.Index)
		userRouter.GET("/userlist", Service.GetUserList)
		userRouter.GET("/adduser", Service.CreateUser)
		userRouter.DELETE("/deluser", Service.DeleteUser)
		userRouter.POST("/updateuser", Service.UpdateUser)
	}

	gitRouter := router.Group("/git")
	{
		gitRouter.GET("/login", Service.GitLogin)
		gitRouter.GET("/callback", Service.GitCallBack)
	}
}
