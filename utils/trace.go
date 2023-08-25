package utils

import (
	"context"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

func StartTrace(serviceName string, version string) {
	uptraceDsn := os.Getenv("UPTRACE_DSN")
	if uptraceDsn == "" {
		log.Fatal("UPTRACE_DSN is required")
	}

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(uptraceDsn),
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(version),
	)
}

func ShutdownTrace(ctx context.Context) {
	err := uptrace.Shutdown(ctx)
	if err != nil {
		log.Fatalf("can't shutdown uptrace: %v", err)
	}
}

func StartSpan(ctx context.Context, serviceName, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)

	return tracer.Start(ctx, spanName)
}
