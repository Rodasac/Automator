package utils

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"log"
)

func StartLogger(debug bool) *otelzap.Logger {
	var zapLogger *zap.Logger
	var err error
	if debug {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return otelzap.New(zapLogger, otelzap.WithMinLevel(zapLogger.Level()))
}

func StartLoggerWithCtx(ctx context.Context, debug bool) otelzap.LoggerWithCtx {
	return StartLogger(debug).Ctx(ctx)
}
