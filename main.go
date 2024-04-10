package main

import (
	"Chat/router"
	"Chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	router.RouterInit()
}
