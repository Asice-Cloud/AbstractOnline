package controller

import (
	"Chat/config"
	"Chat/model"
	"Chat/pkg"
	"Chat/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @router /admin/login [get]
func AdminLogin(ctx *gin.Context) {
	name := ctx.Query("name")
	password := ctx.Query("password")
	if name == "admin" && password == "admin" {
		tokenString, err := pkg.GenerateJWT(name, "admin")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not generate token",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "login success",
		})
		SessionSet(ctx, "admin", tokenString)
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "login failed",
		})
	}
}

// BlockIPRetrieval
// Admin Block IP Retrieval
// @Tags Admin
// @Success	200 {string} json{"code","blockip"}
// @router /admin/retrievalblockip [get]
func BlockIPRetrieval(ctx *gin.Context) {
	// get the blocked IP
	blockIp, err := RetrievalBlockIP(ctx)
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "get blocked ip failed",
		})
		return
	}
	ctx.HTML(200, "blockIp.html", gin.H{
		"blockIp": blockIp,
	})

}

// BlockIPRemove
// Admin Block IP Remove
// @Tags Admin
// @Success	200 {string} json{"code","message"}
// @router /admin/deleteblockip [delete]
func BlockIPRemove(ctx *gin.Context) {
	ip := ctx.Query("ip")
	err := RemoveBlockIP(ctx, ip)
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "remove blocked ip failed",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "remove blocked ip success",
	})

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
// @router /admin/userlist [get]
func GetUserList(context *gin.Context) {
	var data []model.UserBasic
	data = service.GetUserList()
	if len(data) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"message": "No user",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}
