package controller

import "github.com/gin-gonic/gin"

func Before(context *gin.Context) {
	context.HTML(200, "verify.html", nil)
}
