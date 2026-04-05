package config

import (
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(conf *Config) {
	// Ensure log directory exists
	logDir := filepath.Dir(conf.LogFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	// Lumberjack for log rotation
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.LogFile,
		MaxSize:    conf.LogMaxSize, // megabytes
		MaxBackups: conf.LogMaxBackups,
		MaxAge:     conf.LogMaxAge,   // days
		Compress:   conf.LogCompress, // disabled by default
	})

	// Log Level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(conf.LogLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	// Encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Multi-core setup: File (JSON) + Console (Color/Console)
	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), w, level),
	}

	// Add console output in development
	if conf.AppEnv == "development" {
		consoleEncoder := zap.NewDevelopmentEncoderConfig()
		consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), zapcore.AddSync(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)
	Log = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(Log)
}
