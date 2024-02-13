package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// index page

// @Tags			home page
// @Success		200	{string} welcome
// @Router			/index [get]
func Welcome(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"message": "Welcome",
	})
}
