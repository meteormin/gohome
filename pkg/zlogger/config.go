package zlogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Config is custom logger configuration
type Config struct {
	Name        string        // logger name
	TimeFormat  string        // time format
	FilePath    string        // log file path
	Filename    string        // log file name
	MaxSize     int           // max log file size
	MaxBackups  int           // max log file backups
	MaxAge      int           // max log file age
	Compress    bool          // compress log status
	TimeKey     string        // time key
	TimeZone    string        // time zone
	LogLevel    zapcore.Level // log level
	WithOptions []zap.Option  // zap.Option
}

// defaultConfig is default logger configuration
var defaultConfig = Config{
	Name:       "default",
	TimeFormat: "2006-01-02 15:04:05",
	FilePath:   "",
	Filename:   "",
	MaxSize:    10,
	MaxBackups: 30,
	MaxAge:     1,
	Compress:   false,
	TimeKey:    "timestamp",
	TimeZone:   "",
	LogLevel:   zapcore.DebugLevel,
}

// resolveConfig returns the default configuration with optional overrides.
//
// It takes a variable number of Config parameters and returns a Config type.
func resolveConfig(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}

	cfg := config[0]

	if cfg.TimeFormat == "" {
		cfg.TimeFormat = defaultConfig.TimeFormat
	}

	if cfg.TimeKey == "" {
		cfg.TimeKey = defaultConfig.TimeKey
	}

	if cfg.TimeZone == "" {
		cfg.TimeZone, _ = time.Now().Zone()
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultConfig.MaxSize
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = defaultConfig.MaxAge
	}

	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = defaultConfig.MaxBackups
	}

	return cfg
}
