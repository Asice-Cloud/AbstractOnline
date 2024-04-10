package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"Chat/Model"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	rdb *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("Config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to init, err is %v", err)
	}
	fmt.Println("config app", viper.Get("app"))
	fmt.Println("config mysql", viper.Get("mysql"))
}

func InitMySQL() {
	// customize SQL log
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // slow SQL threshold
			LogLevel:      logger.Info, // level
			Colorful:      true,        // use colorful
		},
	)

	// connect to database
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate schema
	if err := DB.AutoMigrate(&Model.UserBasic{}); err != nil {
		return
	}
	fmt.Println("Successfully init")
}

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
	fmt.Println("Successfully init")
}
