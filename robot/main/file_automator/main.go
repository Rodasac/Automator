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
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

// This is meant to be used for testing purposes only.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	debug := os.Getenv("APP_DEBUG")

	var zapLogger *zap.Logger
	if debug == "true" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	browser := rod.New().Context(ctx)
	err = browser.Connect()
	if err != nil {
		zapLogger.Fatal("error connecting to browser", zap.Error(err))
	}
	zapLogger.Debug("Connected to browser")

	pagePool := rod.NewPagePool(3)

	taskController := taskControllers.NewTaskController(browser, pagePool, db, ctx, zapLogger)
	consumerController := controllerConsumer.NewFileConsumerController(taskController, zapLogger)

	go func() {
		errs := consumerController.ConsumeTasks()
		if errs != nil && len(errs) > 0 {
			zapLogger.Fatal("error processing tasks", zap.Errors("errors", errs))
		}

		// Because this is a file consumer we finish here.
		// But, this may not occur on streams implementations.
		stop()
		pagePool.Cleanup(func(page *rod.Page) {
			err := page.Close()
			if err != nil {
				zapLogger.Error("error closing page", zap.Error(err))
			}
		})
	}()

	select {
	case <-ctx.Done():
		zapLogger.Info("Exiting...")
		stop()

		return
	}
}
