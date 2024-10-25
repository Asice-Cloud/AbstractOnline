package admin_module

import (
	"Abstract/config"
	"Abstract/model"
	"Abstract/response"
	"Abstract/service"
	"Abstract/session"
	"Abstract/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	mu sync.Mutex
)

// AdminLogin
// Admin Login
// @Tags Admin
// @param name query string false "Name"
// @param password query string false "Password"
// @Success	200 {string} json{"code","message"}
// @router /admin_module/login [get]
func AdminLogin(ctx *gin.Context) {
	name := ctx.Query("name")
	password := ctx.Query("password")
	if name == "admin" && password == "admin" {
		atoken, rtoken, err := utils.GenToken(0, "admin")
		if err != nil {
			response.RespErrorWithMsg(ctx, response.CodeInvalidToken, errors.New("could not generate token"))
			return
		}
		admin := session.AdminSession{
			ID:           0,
			AdminName:    name,
			AccessToken:  atoken,
			RefreshToken: rtoken,
		}
		session.SessionSet("admin", ctx, "admin", admin)
		response.RespSuccess(ctx, admin)
	} else {
		response.RespSuccess(ctx, response.CodeInvalidPassword)
	}
}

// BlockIPRetrieval
// Admin Block IP Retrieval
// @Tags Admin
// @Success	200 {string} json{"code","blockip"}
// @router /admin-module/retrievalblockip [get]
func BlockIPRetrieval(ctx *gin.Context) {
	// get the blocked IP
	blockIp, err := RetrievalBlockIP(ctx)
	if err != nil {
		response.RespErrorWithMsg(ctx, response.CodeServerBusy, errors.New("get blocked ip failed"))
		return
	}
	response.RespSuccess(ctx, blockIp)

}

// BlockIPRemove
// Admin Block IP Remove
// @Tags Admin
// @Success	200 {string} json{"code","message"}
// @router /admin_module/deleteblockip [delete]
func BlockIPRemove(ctx *gin.Context) {
	ip := ctx.Query("ip")
	err := RemoveBlockIP(ctx, ip)
	if err != nil {
		response.RespErrorWithMsg(ctx, response.CodeServerBusy, errors.New("remove blocked ip failed"))
		return
	}
	response.RespSuccess(ctx, "remove blocked ip success")
}

func RetrievalBlockIP(ctx *gin.Context) (map[string]string, error) {
	mu.Lock()
	defer mu.Unlock()
	return config.Rdb.HGetAll(ctx, "blockip").Result()
}

// clear the blocked IP by hands
func RemoveBlockIP(ctx *gin.Context, ip string) error {
	mu.Lock()
	defer mu.Unlock()
	return config.Rdb.Del(ctx, ip).Err()
}

// GetUserList
// @Summary Find all users
// @Tags Admin
// @Success	200	{string} json{"code","message"}
// @router /admin_module/userlist [get]
func GetUserList(context *gin.Context) {
	var data []model.UserBasic
	data = service.GetUserList()
	if len(data) == 0 {
		response.RespError(context, response.CodeUserNotExist)
		return
	}
	response.RespSuccess(context, data)
}
