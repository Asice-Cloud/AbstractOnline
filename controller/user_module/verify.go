package user_module

import "github.com/gin-gonic/gin"

// @Tags	home page
// @Success	200	{string} welcome
// @router	/user_module/before [get]
func Before(context *gin.Context) {
	context.HTML(200, "verify.html", nil)
}
