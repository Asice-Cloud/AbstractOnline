package main

import (
	"Chat/Router"
	"Chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	Router.RouterInit()
}
