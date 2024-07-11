package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

var Lg *zap.Logger

func GetTime() string {
	now := time.Now()
	date := fmt.Sprintf("%02d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
	return date
}

// InitLogger initializes the logger
func InitLogger() {
	// Initialize viper to read app.yml
	viper.SetConfigName("app") // name of config file (without extension)
	viper.SetConfigType("yml") // type of the config file
	viper.AddConfigPath(".")   // path to look for the config file in

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	// Retrieve the log path from the config
	logPath := viper.GetString("log.path")
	// Ensure the log directory exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			return
		} // Create the directory with appropriate permissions
	}

	// Correct the log path and ensure it's valid
	correctedLogPath := filepath.Join(logPath, fmt.Sprintf("log_%s.log", GetTime()))
	cfg := zap.NewDevelopmentConfig()
	// Use the retrieved log path and append the log file name with mode
	cfg.OutputPaths = []string{correctedLogPath}

	// Initialize the logger with the configuration
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("Fatal error initializing logger: %w \n", err))
	}

	zap.ReplaceGlobals(logger) // Replace the global logger
}
