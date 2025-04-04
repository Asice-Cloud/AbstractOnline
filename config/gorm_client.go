package config

import (
	"Abstract/model"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os/exec"

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
	for {
		DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
		if err != nil {
			cmd := exec.Command("mysql", "-u", "root", "-proot", "-e", "CREATE DATABASE IF NOT EXISTS usertest;")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				Lg.Error("failed to create db", zap.Error(err))
			}
		} else {
			break
		}
	}
	// migrate schema
	if err := DB.AutoMigrate(&model.UserBasic{}); err != nil {
		return
	}
	fmt.Println("Database successfully init")
}
