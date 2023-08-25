package main

import (
	controllerConsumer "automator-go/robot/adapters/controllers/consumer"
	taskControllers "automator-go/robot/adapters/controllers/tasks"
	utils2 "automator-go/utils"
	"context"
	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"strconv"
)

const serviceName = "robot-stream-automator"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	debugEnv := os.Getenv("APP_DEBUG")
	debug := debugEnv == "true"
	version := os.Getenv("APP_VERSION")
	pagePoolNumberStr := os.Getenv("PAGE_POOL_SIZE")
	pagePoolNumber, err := strconv.Atoi(pagePoolNumberStr)
	if err != nil {
		log.Fatal("PAGE_POOL_SIZE is required")
	}

	utils2.StartTrace(serviceName, version, debug)
	defer utils2.ShutdownTrace(ctx)

	ctx, span := utils2.StartSpan(ctx, serviceName, "root")

	logWithCtx := utils2.StartLoggerWithCtx(ctx, debug)

	db := utils2.OpenDb()

	go func() {
		browser := rod.New().Context(ctx)
		err = browser.Connect()
		if err != nil {
			logWithCtx.Fatal("error connecting to browser", zap.Error(err))
		}
		logWithCtx.Debug("Connected to browser")

		pagePool := rod.NewPagePool(pagePoolNumber)

		taskController := taskControllers.NewTaskController(browser, pagePool, db, ctx, &logWithCtx)
		consumerController := controllerConsumer.NewRabbitConsumerController(taskController, &logWithCtx, ctx)

		errs := consumerController.ConsumeTasks()
		if len(errs) > 0 {
			logWithCtx.Fatal("error processing tasks", zap.Errors("errors", errs))
		}

		pagePool.Cleanup(func(page *rod.Page) {
			err := page.Close()
			if err != nil {
				logWithCtx.Error("error closing page", zap.Error(err))
			}
		})
	}()

	<-stopSignal
	span.End()
	logWithCtx.Info("Shutting down")
}
