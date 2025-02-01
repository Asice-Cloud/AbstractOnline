package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"message": "Welcome",
	})
}
