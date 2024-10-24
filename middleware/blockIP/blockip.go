package blockIP

import (
	"Abstract/config"
	"context"
	"sync"
	"time"
)

var (
	mu  sync.Mutex
	ctx = context.Background()
)

func AddBlockIP(blockip string) error {
	// add the IP into blocked IP for 30 days
	mu.Lock()
	err := config.Rdb.Set(ctx, blockip, "blocked", 720*time.Hour).Err()
	if err != nil {
		return err
	}
	mu.Unlock()
	return nil
}
