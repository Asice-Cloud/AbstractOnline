package router

import (
	"Chat/docs"
	"Chat/middleware"
	"Chat/middleware/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	// middleware
	router.Use(cors.New(middleware.CorsInit()))
	router.Use(gin.Logger())
	router.Use(log.LoggingMiddleware())
	//set session
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "127.0.0.1:9999"
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	Routers(router)
	err := router.Run(":9999")
	if err != nil {
		return
	}
}
