package auth

import (
	"Chat/config"
	"context"
	"sync"
	"time"
)

var (
	mu  sync.Mutex
	ctx = context.Background()
)

func AddBlockIP(blocklist map[int]string) error {
	// add the IP into blocked IP for 30 days
	for _, v := range blocklist {
		mu.Lock()
		blocklist[len(blocklist)] = v
		err := config.Rdb.Set(ctx, v, "blocked", 720*time.Hour).Err()
		if err != nil {
			return err
		}
		mu.Unlock()
	}
	return nil
}

// clear the blocked IP by hands
func RemoveBlockIP(ip string) error {
	mu.Lock()
	defer mu.Unlock()
	return config.Rdb.Del(ctx, ip).Err()
}
