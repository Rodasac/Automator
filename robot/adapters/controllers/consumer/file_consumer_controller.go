package consumer

import (
	"automator-go/robot/adapters/controllers/tasks"
	adapterConsumer "automator-go/robot/adapters/gateways/consumer"
	"automator-go/robot/usecases/consumer"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type FileConsumerController struct {
	taskController *tasks.TaskController
	logger         *otelzap.LoggerWithCtx
}

func NewFileConsumerController(
	taskController *tasks.TaskController,
	logger *otelzap.LoggerWithCtx,
) FileConsumerController {
	return FileConsumerController{taskController: taskController, logger: logger}
}

func (f FileConsumerController) ConsumeTasks() []error {
	f.logger.Info("starting consumer")
	consumerHandler := adapterConsumer.NewTaskQueueConsumerFromJSONFile(f.taskController, f.logger)
	consumerUseCase := consumer.NewTaskQueueConsumer(consumerHandler)

	return consumerUseCase.StartConsumer()
}
