package rzap

import (
	"fmt"

	"github.com/rucciva/go-kit/log"

	"go.uber.org/zap"
)

var (
	_ log.BaseLogger = baseLogger{}
	_ log.Logger     = Logger{}
)

type baseLogger struct {
	zap    *zap.SugaredLogger
	fields []interface{}
}

func (l baseLogger) Debug(msg string) {
	l.zap.Debugw(msg, l.fields...)
}
func (l baseLogger) Info(msg string) {
	l.zap.Infow(msg, l.fields...)
}
func (l baseLogger) Warn(msg string) {
	l.zap.Warnw(msg, l.fields...)
}
func (l baseLogger) Error(msg string) {
	l.zap.Errorw(msg, l.fields...)
}
func (l baseLogger) Fatal(msg string) {
	l.zap.Fatalw(msg, l.fields...)
}
func (l baseLogger) Panic(msg string) {
	l.zap.Panicw(msg, l.fields...)
}

func (l baseLogger) Debugf(format string, args ...interface{}) {
	l.zap.Debugw(
		fmt.Sprintf(format, args...),
		l.fields...)
}
func (l baseLogger) Infof(format string, args ...interface{}) {
	l.zap.Infow(
		fmt.Sprintf(format, args...),
		l.fields...)
}

func (l baseLogger) Warnf(format string, args ...interface{}) {
	l.zap.Warnw(
		fmt.Sprintf(format, args...),
		l.fields...)
}
func (l baseLogger) Errorf(format string, args ...interface{}) {
	l.zap.Errorw(
		fmt.Sprintf(format, args...),
		l.fields...)
}
func (l baseLogger) Fatalf(format string, args ...interface{}) {
	l.zap.Fatalw(
		fmt.Sprintf(format, args...),
		l.fields...)
}
func (l baseLogger) Panicf(format string, args ...interface{}) {
	l.zap.Panicw(
		fmt.Sprintf(format, args...),
		l.fields...)
}

type Logger struct {
	baseLogger
	cfg loggerConfig
}

func NewLogger(opts ...LoggerOption) *Logger {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	logger := zap.New(cfg.zapCore(), cfg.opts...)
	return &Logger{baseLogger: baseLogger{zap: logger.Sugar()}, cfg: *cfg}
}

func (l Logger) WithFields(KeyValuePairs ...interface{}) log.BaseLogger {
	return baseLogger{
		zap:    l.zap,
		fields: KeyValuePairs,
	}
}

func (l Logger) Close() error {
	return l.zap.Sync()
}

func (l *Logger) Zap() *zap.Logger {
	return l.zap.Desugar()
}
