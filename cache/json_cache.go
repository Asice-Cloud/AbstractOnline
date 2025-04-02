package cache

import (
	"Abstract/config"
	"context"
	"github.com/bytedance/sonic"
	"time"
)

var rdb = config.Rdb

// SetCache, set cache for common pages or queries
func SetJson(key string, value interface{}, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := sonic.Marshal(value)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, val, expiry).Err()
}

// GetCache
func GetJson(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return sonic.Unmarshal([]byte(val), dest)
}
