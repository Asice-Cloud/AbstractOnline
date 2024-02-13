package Controller

import (
	"Chat/Model"
	"Chat/Service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// User Home Page

// @Tags			User Home
// @Success		200	{string} welcome,user
// @Router			/user/index [get]
func Index(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Welcome!",
	})
}

// GetUserList

// @Tags			UserList
// @Success		200	{string} json{"code","message"}
// @Router			/user/userlist [get]
func GetUserList(context *gin.Context) {
	var data []Model.UserBasic
	data = Service.GetUserList()
	if len(data) != 0 {
		context.JSON(http.StatusOK, gin.H{
			"message": data,
		})
	}
}
