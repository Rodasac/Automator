package main

import (
	taskControllers "automator-go/adapters/controllers/tasks"
	"automator-go/adapters/gateways/browser_automator"
	"automator-go/adapters/gateways/consumer"
	"automator-go/adapters/gateways/hasher"
	storage2 "automator-go/adapters/gateways/storage"
	bunRepo "automator-go/adapters/repositories/bun"
	"automator-go/usecases/task"
	"context"
	"database/sql"
	"github.com/go-rod/rod"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := os.Getenv("DATABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	browser := rod.New()
	automator := browser_automator.NewRodAutomator(browser)
	storage := storage2.NewFileStorage()
	mediaRepo := bunRepo.NewBunCaptureMedia(db, &ctx)
	hashHandler := hasher.NewPHashHandler()
	taskUseCase := task.NewProcessor(automator, mediaRepo, storage, hashHandler)
	taskController := taskControllers.NewTaskController(taskUseCase)

	consumerHandler := consumer.NewTaskQueueConsumerFromJSONFile(taskController)
	err = consumerHandler.ConsumeTasks()
	if err != nil {
		log.Fatal(err)
	}
}
