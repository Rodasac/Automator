package utils

import (
	"context"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
)

func StartTrace(serviceName, version string, debug bool) {
	uptraceDsn := os.Getenv("UPTRACE_DSN")
	if uptraceDsn == "" {
		log.Fatal("UPTRACE_DSN is required")
	}

	env := "production"
	if debug {
		env = "development"
	}

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(uptraceDsn),
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(version),
		uptrace.WithResourceAttributes(
			attribute.String("deployment.environment", env),
		),
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
