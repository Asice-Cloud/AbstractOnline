package websocket_work

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

var names = []string{
	"鸢一折纸",
	"本条二亚",
	"时崎狂三",
	"冰芽川四糸乃",
	"五河琴里",
	"星宫六喰",
	"镜野七罪",
	"八舞耶俱矢",
	"八舞夕弦",
	"诱宵美九",
	"夜刀神十香",
	"夜刀神天香",
	"园神凛绪",
	"园神凛祢",
	"万由理",
	"我简直就是五河士道本人",
	"隐藏角色伊藤诚",
}

//var nameCount = make(map[string]int)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Redis address
	Password: "",               // No password
	DB:       0,                // Default DB
})

func getRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

func getUniqueName(baseName string) string {
	ctx := context.Background()
	key := fmt.Sprintf("name:%s", baseName)
	val, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		fmt.Printf("Failed to increment name count in Redis: %v\n", err)
		return baseName // Fallback to baseName if Redis fails
	}
	if val == 1 {
		return baseName
	}
	return fmt.Sprintf("%s Code.%02d", baseName, val)
}
