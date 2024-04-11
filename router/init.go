package router

import (
	"Chat/docs"
	"Chat/middleware"
	"Chat/middleware/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	// middleware
	router.Use(cors.New(middleware.RouterInit()))
	router.Use(gin.Logger())
	router.Use(log.LoggingMiddleware())

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "127.0.0.1:9999"
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	Routers(router)
	router.Run(":9999")
}
