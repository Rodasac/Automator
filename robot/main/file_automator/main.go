package main

import (
	controllerConsumer "automator-go/robot/adapters/controllers/consumer"
	taskControllers "automator-go/robot/adapters/controllers/tasks"
	"context"
	"database/sql"
	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func shutDownUptrace(ctx context.Context) {
	err := uptrace.Shutdown(ctx)
	if err != nil {
		log.Fatalf("can't shutdown uptrace: %v", err)
	}
}

// This is meant to be used for testing purposes only.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	debug := os.Getenv("APP_DEBUG")

	uptraceDsn := os.Getenv("UPTRACE_DSN")
	if uptraceDsn == "" {
		log.Fatal("UPTRACE_DSN is required")
	}

	version := os.Getenv("APP_VERSION")

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(uptraceDsn),
		uptrace.WithServiceName("robot-file-automator"),
		uptrace.WithServiceVersion(version),
	)

	tracer := otel.Tracer("robot-file-automator")
	ctx, span := tracer.Start(ctx, "root")

	var zapLogger *zap.Logger
	if debug == "true" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	otelLog := otelzap.New(zapLogger, otelzap.WithMinLevel(zapLogger.Level()))
	logWithCtx := otelLog.Ctx(ctx)

	dsn := os.Getenv("DATABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))
	db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("robot")))

	browser := rod.New().Context(ctx)
	err = browser.Connect()
	if err != nil {
		logWithCtx.Fatal("error connecting to browser", zap.Error(err))
	}
	logWithCtx.Debug("Connected to browser")

	pagePool := rod.NewPagePool(3)

	taskController := taskControllers.NewTaskController(browser, pagePool, db, ctx, &logWithCtx)
	consumerController := controllerConsumer.NewFileConsumerController(taskController, &logWithCtx)

	go func() {
		errs := consumerController.ConsumeTasks()
		if errs != nil && len(errs) > 0 {
			logWithCtx.Fatal("error processing tasks", zap.Errors("errors", errs))
		}

		// Because this is a file consumer we finish here.
		// But, this may not occur on streams implementations.
		pagePool.Cleanup(func(page *rod.Page) {
			err := page.Close()
			if err != nil {
				logWithCtx.Error("error closing page", zap.Error(err))
			}
		})
		span.End()
		shutDownUptrace(ctx)
		stop()
	}()

	select {
	case <-ctx.Done():
		logWithCtx.Info("Exiting...")

		return
	}
}
