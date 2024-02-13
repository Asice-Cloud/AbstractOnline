package main

import (
	"Chat/Router"
	"Chat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	Router.RouterInit()
}
