package websocket_work

import (
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
}

func getRandomName() string {
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}
