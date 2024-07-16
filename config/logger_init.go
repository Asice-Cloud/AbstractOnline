package config

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
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
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

func loadConfig(filename string) (*Config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read log config: %v", err)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
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

	var core zapcore.Core
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(config.LogOutput, config.LogFile),
		MaxBackups: config.MaxBackups,
		MaxSize:    config.LogMaxSize,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	if gin.Mode() == gin.DebugMode {
		consoleWriter := zapcore.Lock(os.Stdout)
		multiWriter := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)
		core = zapcore.NewCore(encoder, multiWriter, zap.NewAtomicLevelAt(parseLogLevel(config.LogLevel)))
	} else {
		core = zapcore.NewCore(encoder, fileWriter, zap.NewAtomicLevelAt(parseLogLevel(config.LogLevel)))
	}

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	Lg = logger

	// Create a Pipe to capture stderr
	readEnd, writeEnd, err := os.Pipe()
	if err != nil {
		panic(fmt.Sprintf("Failed to create pipe: %v", err))
	}

	// Redirect stderr to the write end of the pipe
	os.Stderr = writeEnd
	// Capture stderr output in a separate goroutine
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := readEnd.Read(buffer)
			if err != nil {
				if err != io.EOF {
					Lg.Error("Error reading from stderr pipe", zap.Error(err))
				}
				break
			}
			if n > 0 {
				Lg.Error("Stderr", zap.ByteString("output", buffer[:n]))
			}
		}
	}()

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
