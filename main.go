package main

import (
	"Chat/config"
	"Chat/router"
)

func main() {
	config.InitConfig()
	config.InitMySQL()
	config.InitRedis()
	router.RouterInit()
}
