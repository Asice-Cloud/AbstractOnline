package router

import (
	"Abstract/docs"
	"Abstract/middleware/auth"
	"Abstract/middleware/log"
	"Abstract/middleware/safe"
	"Abstract/session"
	"Abstract/utils"
	"fmt"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func RouterInit() {
	router := gin.Default()

	// middleware
	router.Use(cors.New(auth.CorsInit()))
	//router.Use(cors.Default())
	router.Use(gin.Logger())

	//set session
	session.InitSession(router)
	router.Use(log.GinLogger(), log.GinRecovery(true))
	router.Use(safe.SetCSRFToken())
	router.Use(safe.SanitizeInputMiddleware())
	addr := "0.0.0.0:8088"
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = addr
	docs.SwaggerInfo.BasePath = ""

	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	Routers(router)

	utils.Try(func() {
		err := router.Run(addr)
		if err != nil {
			utils.Throw(err)
		}
	}).CatchAll(func(err error) {
		fmt.Printf("Caught: %v\n", err)
	}).Finally(func() {
		fmt.Println("finally")
	})
}
