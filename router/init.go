package router

import (
	"Chat/docs"
	"Chat/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	// middleware
	router.Use(cors.New(middleware.RouterInit()))
	router.Use(gin.Logger())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.LimitCount)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "127.0.0.1:9999"
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("template/*")
	router.Static("/assert", "./assert")

	Routers(router)
	router.Run(":9999")
}
