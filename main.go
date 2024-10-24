package main

import (
	"Abstract/config"
	"Abstract/router"
)

func main() {
	config.InitMode()
	router.RouterInit()
}
