package config

import (
	"Abstract/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
	"time"
)

func initMySQL() {
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
	dsn := "root:123456@tcp(127.0.0.1:3306)/usertest?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	// migrate schema
	if err := DB.AutoMigrate(&model.UserBasic{}); err != nil {
		return
	}
	fmt.Println("Database successfully init")
}
