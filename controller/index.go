package Service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// index page

// @Tags	home page
// @Success	200	{string} welcome
// @router	/index [get]
func Welcome(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"message": "Welcome",
	})
}
