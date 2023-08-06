package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements the Logger interface using zap logging package.
type zapLogger struct {
	logger *zap.SugaredLogger
}

var _ Logger = (*zapLogger)(nil)

// New - creates new instance logger.
func New(level string) Logger {
	var l zapcore.Level
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		l = zap.InfoLevel
	}

	// logger config
	config := zap.Config{
		Development:      false,
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder, // e.g. "Info"
			CallerKey:      "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder, // e.g. package/file:line
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder, // e.g. 2020-05-05T03:24:36.903+0300
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder, // e.g. GetSiteGeneralHandler
			MessageKey:     "message",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	// build logger from config
	logger, _ := config.Build()

	// configure and create logger
	return &zapLogger{
		logger: logger.Sugar(),
	}
}

// Named - returns a new logger with a chained name.
func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{
		logger: l.logger.Named(name),
	}
}

// With - returns a new logger with parameters.
func (l *zapLogger) With(args ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(args...),
	}
}

// WithContext - returns a new logger with a chained name.
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	return l.With("RequestID", ctx.Value("RequestID"))
}

// Debug - logs in debug level.
func (l *zapLogger) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
}

// Info - logs in info level.
func (l *zapLogger) Info(message string, args ...interface{}) {
	l.logger.Infow(message, args...)
}

// Warn - logs in warn level.
func (l *zapLogger) Warn(message string, args ...interface{}) {
	l.logger.Warnw(message, args...)
}

// Error - logs in error level.
func (l *zapLogger) Error(message string, args ...interface{}) {
	l.logger.Errorw(message, args...)
}

// Fatal - logs and exits program with status 1.
func (l *zapLogger) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
	os.Exit(1)
}

func (l *zapLogger) Unwrap() *zap.Logger {
	return l.logger.Desugar()
}
