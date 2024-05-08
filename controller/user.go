package controller

import (
	"Chat/model"
	"Chat/response"
	"Chat/service"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserSession struct {
	UserID       int
	UserName     string
	AccessToken  string
	RefreshToken string
}

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
	data, err := service.Login(name, password)
	if err != nil {
		response.RespError(context, response.CodeInvalidPassword)
		return
	}
	user, ok := data.(model.UserBasic)
	if !ok {
		response.RespError(context, response.CodeInvalidToken)
	}
	// Set the session for the user
	userSession := UserSession{
		UserID:       int(user.ID),
		UserName:     user.Name,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
	SessionSet(context, fmt.Sprintf("user_%d", user.ID), userSession)
	response.RespSuccess(context, userSession)
}

// Logout
// @Summary	user login
// @Tags UserModule
// @param userID query string false "userID"
// @Success	200	{string} json{"code","message"}
// @router /user/logout [delete]
func Logout(context *gin.Context) {
	userID, _ := strconv.Atoi(context.Query("userID"))
	// Get the session for the user
	session := sessions.Default(context)
	userSession := session.Get(fmt.Sprintf("user_%d", userID))
	if userSession == nil {
		response.RespError(context, response.CodeNotLogin)
		return
	}
	// Delete the session for the user
	session.Delete(fmt.Sprintf("user_%d", userID))
	session.Save()
	response.RespSuccess(context, "Successfully logout")
}

// Index
// User Home Page
// @Tags User Home
// @Success	200	{string} welcome, user
// @router /user/index [get]
func Index(context *gin.Context) {
	response.RespSuccess(context, "Welcome")
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
		response.RespErrorWithMsg(context, response.CodeInvalidParams, errors.New("twice passwords are not consistent"))
		return
	}
	user.Password = password
	rep, err := service.CreatUser(user)
	if rep == -1 {
		response.RespError(context, response.CodeUserNotExist)
		return
	}
	if err != nil {
		response.RespErrorWithMsg(context, response.CodeInvalidParams, err)
		return
	}
	response.RespSuccess(context, fmt.Sprintf("Successfully create user,ID: %v", &user.ID))
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
		response.RespErrorWithMsg(context, response.CodeInvalidParams, errors.New("please input a valid number"))
		return
	}
	user.ID = uint(id)
	err = service.DeleteUser(user)
	if err != nil {
		response.RespError(context, response.CodeServerBusy)
		return
	}
	response.RespSuccess(context, "Successfully delete user")
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

	userSession := SessionGet(context, fmt.Sprintf("user_%d", user.ID))
	if userSession == nil {
		response.RespError(context, response.CodeNotLogin)
		return
	}

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		response.RespError(context, response.CodeInvalidParams)
		return
	} else {
		rep, err := service.UpdateUser(user)
		if err != nil {
			response.RespErrorWithMsg(context, response.CodeServerBusy, errors.New("failed to update user"))
			return
		}
		if rep == -1 {
			response.RespError(context, response.CodeUserNotExist)
		}
		// Update the session for the user
		userData := rep.(model.UserBasic)
		data := UserSession{
			UserID:       int(userData.ID),
			UserName:     userData.Name,
			AccessToken:  userData.AccessToken,
			RefreshToken: userData.RefreshToken,
		}
		SessionUpdate(context, fmt.Sprintf("user_%d", user.ID), data)
		response.RespSuccess(context, fmt.Sprintf("Successfully update user,ID: %d", data))
	}
}

/**
 * @Author huchao
 * @Description //TODO 刷新accessToken
 * @Date 17:09 2022/2/17
 **/
// RefreshTokenHandler 刷新accessToken
// @Summary 刷新accessToken
// @Description 刷新accessToken
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /refresh_token [GET]
/*func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
*/
