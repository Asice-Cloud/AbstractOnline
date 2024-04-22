package controller

import (
	"Chat/config"
	"github.com/gin-gonic/gin"
	"sync"
)

var (
	mu sync.Mutex
)

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
