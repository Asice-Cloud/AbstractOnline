package admin_module

import (
	"Abstract/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
)

var (
	mu sync.Mutex
)

func BlockIPRetrieval(ctx *gin.Context) {
	// get the blocked IP
	blockIp, err := RetrievalBlockIP(ctx)
	if err != nil {
		config.Lg.Debug("retrieval blocked ip failed", zap.Error(err))
		return
	}
	config.Lg.Debug("retrieval blocked ip", zap.Any("blockip", blockIp))
}

func BlockIPRemove(ctx *gin.Context) {
	ip := ctx.Query("ip")
	err := RemoveBlockIP(ctx, ip)
	if err != nil {
		config.Lg.Debug("remove blocked ip failed", zap.Error(err))
		return
	}
	config.Lg.Debug("remove blocked ip success", zap.String("ip", ip))
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
