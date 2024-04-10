package config

import (
	"Chat/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
	"time"
)

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
	if err := DB.AutoMigrate(&model.UserBasic{}); err != nil {
		return
	}
	fmt.Println("Database successfully init")
}
