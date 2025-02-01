package router

import (
	"Abstract/middleware/auth"
	logger "Abstract/middleware/log"
	"Abstract/middleware/safe"
	"Abstract/session"
	"Abstract/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	router.Use(logger.GinLogger(), logger.GinRecovery(true))
	router.Use(safe.SetCSRFToken())
	router.Use(safe.SanitizeInputMiddleware())
	router.LoadHTMLGlob("template/*")
	router.Static("/static", "./static")

	Routers(router)

	utils.Try(func() {
		srv := http.Server{
			Addr:    ":9999",
			Handler: router,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				utils.Throw(err)
			}
		}()
		// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}).CatchAll(func(err error) {
		log.Fatalf("Caught: %v\n", err)
	}).Finally(func() {
		fmt.Println("finally")
	})
}
