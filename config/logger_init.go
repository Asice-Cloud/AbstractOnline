package config

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
)

var Lg *zap.Logger

type Config struct {
	LogLevel   string `json:"log_level"`
	LogOutput  string `json:"log_output"`
	LogFile    string `json:"log_file"`
	LogMaxSize int    `json:"log_max_size"`
	TimeFormat string `json:"time_format"`
	LogFormat  string `json:"log_format"`
	Compress   bool   `json:"compress"`
}

func loadConfig(filename string) (*Config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read log config: %v", err)
		return nil, err
	}

	var config Config
	json.Unmarshal(configFile, &config)
	return &config, nil
}

func InitLogger() {
	config, err := loadConfig("log.json")
	if err != nil {
		panic(fmt.Sprintf("Could not load config: %v", err))
	}

	// Ensure the log directory exists
	if _, err := os.Stat(config.LogOutput); os.IsNotExist(err) {
		err := os.MkdirAll(config.LogOutput, 0755)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log directory: %v", err))
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	// Convert the TimeFormat from log.json to Go's time format
	goTimeFormat := "2006-01-02 15:04:05" // This should match your desired format
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(goTimeFormat)

	var encoder zapcore.Encoder
	if config.LogFormat == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(&lumberjack.Logger{
			Filename: filepath.Join(config.LogOutput, config.LogFile),
			MaxSize:  config.LogMaxSize, // megabytes
			Compress: config.Compress,
		}),
		zap.NewAtomicLevelAt(parseLogLevel(config.LogLevel)),
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
	Lg = logger
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
