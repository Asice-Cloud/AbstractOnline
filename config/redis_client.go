package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() {
	// init redis config:
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
			PoolSize: 100,
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("failed to connect redis")
	}
	fmt.Println("Redis successfully init")
}
