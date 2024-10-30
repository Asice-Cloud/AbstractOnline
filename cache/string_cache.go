package cache

import (
	"Abstract/config"
	"context"
	"go.uber.org/zap"
	"time"
)

func SetString(ctx context.Context, key, value string) {
	err := rdb.Set(ctx, key, value, 1800*time.Second).Err()
	if err != nil {
		config.Lg.Error("could not set cache", zap.String("reason", err.Error()))
		return
	}
}

func GetString(ctx context.Context, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		config.Lg.Error("could not get cache", zap.String("reason", err.Error()))
		return "", err
	}
	return val, nil
}

func GetAllString(ctx context.Context, key string) []string {
	keys := []string{}
	iter := rdb.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		config.Lg.Error("could not get cache", zap.String("reason", err.Error()))
		return nil
	}
	return keys
}
