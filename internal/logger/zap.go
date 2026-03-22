package logger

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	levelKey   = "level"
	messageKey = "message"
	timeKey    = "time"
)

var zapLogLevelMap = map[LogLevel]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	PanicLevel: zapcore.PanicLevel,
	FatalLevel: zapcore.FatalLevel,
}

func NewZapLogger(level LogLevel, w io.Writer) (Logger, error) {
	zapLevel, ok := zapLogLevelMap[level]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %v", level)
	}

	zapLogger, err := initZapLogger(zapLevel, w)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zap logger: %w", err)
	}

	sugarLogger := zapLogger.Sugar()

	return &zapLoggerImpl{
		logger: sugarLogger,
	}, nil
}

type zapLoggerImpl struct {
	logger *zap.SugaredLogger
}

func (l *zapLoggerImpl) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l *zapLoggerImpl) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l *zapLoggerImpl) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}
func (l *zapLoggerImpl) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l *zapLoggerImpl) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
func (l *zapLoggerImpl) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *zapLoggerImpl) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}
func (l *zapLoggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}
func (l *zapLoggerImpl) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}
func (l *zapLoggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}
func (l *zapLoggerImpl) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}
func (l *zapLoggerImpl) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func initZapLogger(level zapcore.Level, w io.Writer) (*zap.Logger, error) {
	writeSyncer := zapcore.AddSync(w)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = levelKey
	encoderConfig.MessageKey = messageKey
	encoderConfig.TimeKey = timeKey
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, level)

	logger := zap.New(core)

	return logger, nil
}
