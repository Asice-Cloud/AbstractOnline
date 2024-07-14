package main

import (
	"Chat/config"
	"Chat/router"
)

func main() {
	config.InitMode()
	config.InitConfig()
	config.InitMySQL()
	config.InitRedis()
	config.InitLogger()
	router.RouterInit()
}
