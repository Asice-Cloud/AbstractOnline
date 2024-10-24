package main

import (
	"Abstract/config"
	"Abstract/router"
)

func main() {
	config.InitMode()
	config.InitConfig()
	config.InitMySQL()
	config.InitRedis()
	config.InitLogger()
	router.RouterInit()
}
