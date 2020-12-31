package rzap

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/rucciva/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerConfig struct {
	enc zapcore.Encoder
	ws  zapcore.WriteSyncer
	lvl zapcore.LevelEnabler

	opts []zap.Option
}

func defaultConfig() *loggerConfig {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339Nano))
	}

	return &loggerConfig{
		lvl: zap.NewAtomicLevel(),
		ws:  zapcore.Lock(os.Stdout),
		enc: zapcore.NewJSONEncoder(encCfg),
	}

}

func (cfg *loggerConfig) zapCore() zapcore.Core {
	return zapcore.NewCore(cfg.enc, cfg.ws, cfg.lvl)
}

func (cfg *loggerConfig) zapOptions() []zap.Option {
	return cfg.opts
}

type LoggerOption func(cfg *loggerConfig)

func WithLevel(l log.Level) LoggerOption {
	return func(cfg *loggerConfig) {
		switch l {
		case log.DebugLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case log.InfoLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		case log.WarnLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		case log.ErrorLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		case log.FatalLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		case log.PanicLevel:
			cfg.lvl = zap.NewAtomicLevelAt(zapcore.PanicLevel)
		}
	}
}

func WithIODiscard() LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.ws = zapcore.AddSync(ioutil.Discard)
	}
}

// WithTWritter use `tWritter` which use `*testing.T`'s `Log()`
func WithTWritter(t *testing.T) LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.ws = &tWritter{t}
	}
}

// WithCapturer set `WriteSyncer` that can capture each line writen to it and optionally forward the line to previously setted `WriteSyncer`
func WithCapturer(forward bool) LoggerOption {
	return func(cfg *loggerConfig) {
		if forward {
			cfg.ws = newCapturer(cfg.ws)
			return
		}
		cfg.ws = newCapturer(nil)
	}
}

func WithConsoleEncoder() LoggerOption {
	return func(cfg *loggerConfig) {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339Nano))
		}
		cfg.enc = zapcore.NewConsoleEncoder(encoderCfg)
	}
}
func WithErrorStackTrace() LoggerOption {
	return func(cfg *loggerConfig) {
		if cfg.enc == nil {
			WithConsoleEncoder()(cfg)
		}
		cfg.enc = errorVerboseToStacktraceEncoder{cfg.enc}
	}
}

func WithCallerInfo() LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.opts = append(cfg.opts, zap.AddCaller(), zap.AddCallerSkip(1))
	}
}

func WithCallerSkip(skip int) LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.opts = append(cfg.opts, zap.AddCallerSkip(skip))
	}
}
