package utils

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"log"
)

func StartLogger(ctx context.Context, debug string) otelzap.LoggerWithCtx {
	var zapLogger *zap.Logger
	var err error
	if debug == "true" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	otelLog := otelzap.New(zapLogger, otelzap.WithMinLevel(zapLogger.Level()))

	return otelLog.Ctx(ctx)
}
