package router

import (
	"Chat/controller"
	"Chat/middleware/blockIP"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
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

		userRouter.GET("/ws", controller.WebSocketHandler)
	}

	v1 := router.Group("/v1")
	{
		v1.GET("/login", controller.GitLogin)
		v1.GET("/callback", controller.GitCallBack)
		v1.GET("/error", func(context *gin.Context) {
			context.JSON(500, gin.H{
				"message": "error",
			})
			fmt.Fprintf(os.Stderr, "error")
		})
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
