package zlogger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Default is default logger name
const Default string = "default"

// gLoggers is global logger map
var gLoggers = make(map[string]*zap.SugaredLogger)

// GetLogger returns a *zap.SugaredLogger.
//
// It takes an optional loggerName string parameter and returns a *zap.SugaredLogger.
func GetLogger(loggerName ...string) *zap.SugaredLogger {
	var logger *zap.SugaredLogger
	if gLoggers == nil {
		return New()
	}

	if len(loggerName) == 0 {
		if gLoggers[Default] != nil {
			logger = gLoggers[Default]
		} else {
			for _, l := range gLoggers {
				logger = l
				break
			}
		}
	} else {
		if _, ok := gLoggers[loggerName[0]]; ok {
			logger = gLoggers[loggerName[0]].Named(loggerName[0])
		}
	}

	if logger == nil {
		logger = New()
	}

	return logger
}

// New initializes and returns a new sugared logger.
//
// It accepts a variable number of Config options.
// Returns a pointer to a zap.SugaredLogger.
func New(config ...Config) *zap.SugaredLogger {
	cfg := resolveConfig(config...)

	logFilename := path.Join(cfg.FilePath, cfg.Filename)

	ll := &lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	ws := zapcore.AddSync(ll)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = cfg.TimeKey
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		loc, err := time.LoadLocation(cfg.TimeZone)
		if err != nil {
			loc = time.Local
		}

		t = t.In(loc)

		type appendTimeEncoder interface {
			// AppendTimeLayout description of the Go function.
			//
			// This function takes a time.Time and a string as parameters.
			AppendTimeLayout(time.Time, string)
		}

		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, cfg.TimeFormat)
			return
		}

		enc.AppendString(t.Format(cfg.TimeFormat))
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, ws, cfg.LogLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), cfg.LogLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	logger := zapLogger.Named(cfg.Name).Sugar()

	if cfg.WithOptions != nil {
		logger = logger.WithOptions(cfg.WithOptions...)
	}

	gLoggers[cfg.Name] = logger

	return logger
}
