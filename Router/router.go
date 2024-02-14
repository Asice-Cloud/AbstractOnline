package Router

import (
	"Chat/Controller"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers(router *gin.Engine) {
	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/index", Controller.Welcome)

	userRouter := router.Group("/user")
	{
		userRouter.GET("/index", Controller.Index)
		userRouter.GET("/userlist", Controller.GetUserList)
		userRouter.GET("/adduser", Controller.CreateUser)
		userRouter.DELETE("/deluser", Controller.DeleteUser)
		userRouter.POST("/updateuser", Controller.UpdateUser)
	}

	gitRouter := router.Group("/git")
	{
		gitRouter.GET("/login", Controller.GitLogin)
		gitRouter.GET("/callback", Controller.GitCallBack)
	}
}
