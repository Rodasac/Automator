package consumer

import (
	"automator-go/robot/adapters/controllers/tasks"
	adapterConsumer "automator-go/robot/adapters/gateways/consumer"
	"automator-go/robot/usecases/consumer"
	"automator-go/utils"
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"os"
)

type RabbitConsumerController struct {
	taskController *tasks.TaskController
	logger         *otelzap.LoggerWithCtx
	ctx            context.Context
}

func NewRabbitConsumerController(
	taskController *tasks.TaskController,
	logger *otelzap.LoggerWithCtx,
	ctx context.Context,
) RabbitConsumerController {
	return RabbitConsumerController{taskController: taskController, logger: logger, ctx: ctx}
}

func (r RabbitConsumerController) ConsumeTasks() []error {
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	consumerName := os.Getenv("RABBITMQ_CONSUMER_NAME")
	connectionName := os.Getenv("RABBITMQ_CONNECTION_NAME")
	if queueName == "" || consumerName == "" || connectionName == "" {
		r.logger.Fatal("RABBITMQ_QUEUE_NAME, RABBITMQ_CONSUMER_NAME and RABBITMQ_CONNECTION_NAME are required")
	}

	r.logger.Info("starting consumer")
	c, err := utils.StartClient(r.logger, connectionName)
	if err != nil {
		r.logger.Fatal("Error starting consumer", zap.Error(err))
	}
	defer func(consumer *utils.Consumer) {
		err := consumer.Shutdown()
		if err != nil {
			r.logger.Fatal("Error shutting down consumer", zap.Error(err))
		}
	}(c)

	consumerHandler := adapterConsumer.NewRabbitTaskQueueConsumer(
		c.Channel,
		r.taskController,
		r.logger,
		r.ctx,
		queueName,
		consumerName,
	)
	consumerUseCase := consumer.NewTaskQueueConsumer(consumerHandler)

	return consumerUseCase.StartConsumer()
}
