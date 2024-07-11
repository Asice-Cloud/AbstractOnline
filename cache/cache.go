package cache

import (
	"Chat/config"
	"context"
	"encoding/json"
	"time"
)

// SetCache, set cache for common pages or queries
func SetCache(key string, value interface{}, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return config.Rdb.Set(ctx, key, val, expiry).Err()
}

// GetCache
func GetCache(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := config.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}
