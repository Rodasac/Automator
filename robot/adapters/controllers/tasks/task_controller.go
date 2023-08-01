package tasks

import (
	"automator-go/robot/adapters/gateways/browser_automator"
	"automator-go/robot/adapters/gateways/hasher"
	"automator-go/robot/adapters/gateways/storage"
	bunRepo "automator-go/robot/adapters/repositories/bun"
	"automator-go/robot/entities/models"
	"automator-go/robot/usecases/task"
	"context"
	"github.com/go-rod/rod"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type TaskController struct {
	browser  *rod.Browser
	pagePool rod.PagePool
	db       *bun.DB
	ctx      context.Context
	logger   *zap.Logger
}

func NewTaskController(
	browser *rod.Browser,
	pagePool rod.PagePool,
	db *bun.DB,
	ctx context.Context,
	logger *zap.Logger,
) *TaskController {
	return &TaskController{
		browser:  browser,
		pagePool: pagePool,
		db:       db,
		ctx:      ctx,
		logger:   logger,
	}
}

func (t *TaskController) ProcessTask(taskToProcess *models.Task) error {
	t.logger.Debug("Initializing task processor")
	automator := browser_automator.NewRodAutomator(t.browser, t.pagePool, t.logger)
	fileStorage := storage.NewFileStorage("png", t.logger)
	mediaRepo := bunRepo.NewBunCaptureMedia(t.db)
	hashHandler := hasher.NewPHashHandler(t.logger)
	taskUseCase := task.NewProcessor(automator, mediaRepo, fileStorage, hashHandler)
	t.logger.Debug("Finished initializing task processor")

	return taskUseCase.Process(taskToProcess, t.ctx)
}
