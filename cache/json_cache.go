package cache

import (
	"Abstract/config"
	"context"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var rdb = config.Rdb

// SetCache, set cache for common pages or queries
func SetJson(key string, value interface{}, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := json.Marshal(value)
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

	return json.Unmarshal([]byte(val), dest)
}
