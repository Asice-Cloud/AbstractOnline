package Router

import (
	"Chat/Middleware"
	"Chat/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	//set up CORS Middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                        // allowed origin，use * represent for plural
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // allowed http method
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // allowed http header

	//middleware
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(Middleware.LoggingMiddleware())
	router.Use(Middleware.LimitCount)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "127.0.0.1:9999"
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("Template/*")
	router.Static("/Assert", "./Assert")

	Routers(router)
	router.Run(":9999")
}
