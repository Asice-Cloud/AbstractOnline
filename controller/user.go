package controller

import (
	"Chat/pkg"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"strconv"

	"Chat/model"
	"Chat/service"

	"github.com/gin-gonic/gin"
)

// Login
// @Summary	user login
// @Tags UserModule
// @param name query string false "Name"
// @param password query string false "Password"
// @Success	200	{string} json{"code","message"}
// @router /user/login [get]
func Login(context *gin.Context) {
	name := context.Query("name")
	password := context.Query("password")
	userID, err := service.Login(name, password)
	if err != nil {
		context.JSON(-1, gin.H{
			"message": "Failed to login",
		})
		return
	}
	// Generate JWT token
	token, err := pkg.GenerateJWT(fmt.Sprintf("%d", userID))
	SessionSet(context, "userID", token)
	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully login",
	})
}

// Index
// User Home Page
// @Tags User Home
// @Success	200	{string} welcome, user
// @router /user/index [get]
func Index(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Welcome!",
	})
}

// CreateUser
// @Summary	Add user
// @Tags UserModule
// @param name query string false "Name"
// @param password query string false "Password"
// @param repassword query string false "Twice Password"
// @Success	200	{string} json{"code","message"}
// @router /user/adduser [get]
func CreateUser(context *gin.Context) {
	user := model.UserBasic{}
	user.Name = context.Query("name")
	password := context.Query("password")
	repassword := context.Query("repassword")
	if password != repassword {
		context.JSON(-1, gin.H{
			"message": "twice password is not matched",
		})
		return
	}
	user.Password = password
	rep, err := service.CreatUser(user)
	if rep == -1 {
		context.JSON(-1, gin.H{
			"message": "User already exist",
		})
		return
	}
	if err != nil {
		context.JSON(-1, gin.H{
			"message": "Failed to add user",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully create user: %d", rep),
	})
}

// DeleteUser
// @Summary	Delete user
// @Tags UserModule
// @param id query string false "id"
// @Success	200	{string} json{"code","message"}
// @router /user/deluser [delete]
func DeleteUser(context *gin.Context) {
	var user model.UserBasic
	id, err := strconv.Atoi(context.Query("id"))
	if err != nil {
		context.JSON(-1, gin.H{
			"message": "Please Input a valid number",
		})
		return
	}
	user.ID = uint(id)
	err = service.DeleteUser(user)
	if err != nil {
		context.JSON(-1, gin.H{
			"message": "Failed to delete user",
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
// @Accept x-www-form-urlencoded
// @Param id formData string false "id"
// @Param name formData string false "name"
// @Param password formData string false "password"
// @Param phone formData string false "phone"
// @Param email formData string false "email"
// @Success	200	{string} json{"code","message"}
// @router /user/updateuser [post]
func UpdateUser(context *gin.Context) {
	var user model.UserBasic
	id, _ := strconv.Atoi(context.PostForm("id"))
	user.ID = uint(id)
	user.Name = context.PostForm("name")
	user.Password = context.PostForm("password")
	user.Phone = context.PostForm("phone")
	user.Email = context.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		context.JSON(-1, gin.H{
			"message": "parameters not matched",
		})
		return
	} else {
		rep, err := service.UpdateUser(user)
		if err != nil {
			context.JSON(-1, gin.H{
				"message": "Failed to update user information",
			})
			return
		}
		if rep == -1 {
			context.JSON(-1, gin.H{
				"message": "User not exists",
			})
		}
		context.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully update user: %d", rep),
		})
	}

}
