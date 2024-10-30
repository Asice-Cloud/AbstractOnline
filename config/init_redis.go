package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func initRedis() {
	// init redis config:
	Rdb = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
			PoolSize: 100,
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
}
