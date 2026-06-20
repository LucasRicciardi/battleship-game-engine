package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Logger        *zap.Logger
	SugaredLogger *zap.SugaredLogger
)

// Init initializes the logger with layer-based log levels
func Init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zap.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = "function"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.EncodeLevel = zap.LowercaseLevelEncoder
	config.EncoderConfig.EncodeCaller = zap.ShortCallerEncoder
	config.EncoderConfig.EncodeDuration = zap.StringDurationEncoder
	config.EncoderConfig.EncodeName = zap.FullNameEncoder

	// Set log level based on environment
	if os.Getenv("LOG_LEVEL") != "" {
		level := new(zapcore.Level)
		if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err == nil {
			config.Level.SetLevel(*level)
		}
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	SugaredLogger = Logger.Sugar()
}

// Debug logs a message at Debug level
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info logs a message at Info level
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn logs a message at Warn level
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error logs a message at Error level
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// DPanic logs a message at DPanic level
func DPanic(msg string, fields ...zap.Field) {
	Logger.DPanic(msg, fields...)
}

// Panic logs a message at Panic level
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// Fatal logs a message at Fatal level
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Debugf logs a message at Debug level with formatting
func Debugf(format string, args ...interface{}) {
	SugaredLogger.Debugf(format, args...)
}

// Infof logs a message at Info level with formatting
func Infof(format string, args ...interface{}) {
	SugaredLogger.Infof(format, args...)
}

// Warnf logs a message at Warn level with formatting
func Warnf(format string, args ...interface{}) {
	SugaredLogger.Warnf(format, args...)
}

// Errorf logs a message at Error level with formatting
func Errorf(format string, args ...interface{}) {
	SugaredLogger.Errorf(format, args...)
}

// DPanicf logs a message at DPanic level with formatting
func DPanicf(format string, args ...interface{}) {
	SugaredLogger.DPanicf(format, args...)
}

// Panicf logs a message at Panic level with formatting
func Panicf(format string, args ...interface{}) {
	SugaredLogger.Panicf(format, args...)
}

// Fatalf logs a message at Fatal level with formatting
func Fatalf(format string, args ...interface{}) {
	SugaredLogger.Fatalf(format, args...)
}
