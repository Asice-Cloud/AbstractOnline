package router

import (
	"Chat/docs"
	"Chat/middleware/auth"
	"Chat/middleware/log"
	"Chat/middleware/safe"
	"Chat/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	// middleware
	router.Use(cors.New(auth.CorsInit()))
	router.Use(gin.Logger())
	router.Use(log.LoginMiddleware())
	//set session
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(log.GinLogger(), log.GinRecovery(true))
	router.Use(safe.SetCSRFToken())

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "127.0.0.1:9999"
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	Routers(router)

	utils.Try(func() {
		err := router.Run(":9999")
		if err != nil {
			utils.Throw(err)
		}
	}).CatchAll(func(err error) {
		fmt.Printf("Caught: %v\n", err)
	}).Finally(func() {
		fmt.Println("finally")
	})
}
