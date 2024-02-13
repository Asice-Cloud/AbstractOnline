package Router

import (
	"Chat/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:9999"}                    // 允许的来源，可以是多个域名或 "*" 表示所有来源
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // 允许的 HTTP 方法
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // 允许的 HTTP 标头

	router.Use(cors.New(config))
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
