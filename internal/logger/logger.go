package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Naumovets/backuper/internal/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const CtxKeyLogger = "logger"

func InitLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	var cores []zapcore.Core
	var encoder zapcore.Encoder

	// Choose encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Setup output paths
	for _, path := range cfg.OutputPaths {
		if path == "stdout" {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level))
			continue
		}
		if path == "stderr" {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level))
			continue
		}

		// If rotation is configured, use lumberjack for file outputs
		var w io.Writer
		if cfg.RotationSettings != nil {
			w = &lumberjack.Logger{
				Filename:   path,
				MaxSize:    cfg.RotationSettings.MaxSize,
				MaxBackups: cfg.RotationSettings.MaxCount,
				MaxAge:     cfg.RotationSettings.MaxAge,
				Compress:   cfg.RotationSettings.Compress,
			}
		} else {
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file %s: %w", path, err)
			}
			w = f
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(w), level))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller())

	return logger, nil
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(CtxKeyLogger).(*zap.Logger)
	if !ok {
		logger = zap.L()
		logger.Warn("Logger not found in context", zap.Error(ErrLoggerNotFound))
	}

	return logger
}
