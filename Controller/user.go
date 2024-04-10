package Controller

import (
	"errors"
	"net/http"
	"strconv"

	"Chat/Model"
	"Chat/Service"

	"github.com/gin-gonic/gin"
)

func Login(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
	})
}

// User Home Page
// @Tags User Home
// @Success	200	{string} welcome,user
// @Router /user/index [get]
func Index(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Welcome!",
	})
}

// GetUserList
// @Summary Find all user
// @Tags UserModule
// @Success	 200	{string} json{"code","message"}
// @Router /user/userlist [get]
func GetUserList(context *gin.Context) {
	var data []Model.UserBasic
	data = Service.GetUserList()
	if len(data) != 0 {
		context.JSON(http.StatusOK, gin.H{
			"message": data,
		})
	}
}

// CreateUser
// @Summary	Add user
// @Tags UserModule
// @param name query string false "Name"
// @param password query string false "Password"
// @param repassword query string false "Twice Password"
// @Success	200	{string} json{"code","message"}
// @Router /user/adduser [get]
func CreateUser(context *gin.Context) {
	user := Model.UserBasic{}
	user.Name = context.Query("name")
	password := context.Query("password")
	repassword := context.Query("repassword")
	if password != repassword {
		context.JSON(-1, gin.H{
			"message": errors.New("twice password is not matched"),
		})
		return
	}
	user.Password = password
	rep, err := Service.CreatUser(user)
	if err != nil {
		context.JSON(-1, gin.H{
			"message": errors.New("Failed to add user"),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Add user: " + rep.(string) + " successfully",
	})
}

// DeleteUser
// @Summary	Delete user
// @Tags UserModule
// @param id query string false "id"
// @Success	200	{string} json{"code","message"}
// @Router /user/deluser [delete]
func DeleteUser(context *gin.Context) {
	var user Model.UserBasic
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		context.JSON(-1, gin.H{
			"message": errors.New("Please Input a valid number"),
		})
		return
	}
	user.ID = uint(id)
	err = Service.DeleteUser(user)
	if err != nil {
		context.JSON(-1, gin.H{
			"message": errors.New("Failed to delete user"),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully delete user",
	})
}

// UpdateUser
// @Summary	Update user
// @Tags UserModule
// @param id query string false "id"
// @param name query string false "name"
// @param password query string false "password"
// @Success	200	{string} json{"code","message"}
// @Router /user/updateuser [post]
func UpdateUser(context *gin.Context) {
	var user Model.UserBasic
	id, err := strconv.Atoi(context.PostForm("id"))
	if err != nil {
		context.JSON(-1, gin.H{
			"message": errors.New("Please input valid number"),
		})
		return
	}
	user.ID = uint(id)
	user.Name = context.PostForm("name")
	user.Password = context.PostForm("password")
	rep, err := Service.UpdateUser(user)
	if err != nil {
		context.JSON(-1, gin.H{
			"message": errors.New("Failed to update user information"),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully update user: " + rep.(string),
	})
}
