package logger

import (
	"backend/core/config"

	"go.uber.org/zap"
)

func New(cfg *config.Config) *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment(zap.WithCaller(true), zap.AddStacktrace(zap.ErrorLevel))

	if cfg.LogLevel != "debug" {
		logger, _ = zap.NewProduction(zap.WithCaller(true), zap.AddStacktrace(zap.ErrorLevel))
	}

	return logger.Sugar()
}
