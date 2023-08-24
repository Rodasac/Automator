package consumer

import (
	"automator-go/robot/adapters/controllers/tasks"
	adapterConsumer "automator-go/robot/adapters/gateways/consumer"
	"automator-go/robot/main/utils"
	"automator-go/robot/usecases/consumer"
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
	bindingKey := os.Getenv("RABBITMQ_BINDING_KEY")
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	if queueName == "" || consumerName == "" || bindingKey == "" || exchange == "" {
		r.logger.Fatal("RABBITMQ_QUEUE_NAME, RABBITMQ_CONSUMER_NAME, RABBITMQ_BINDING_KEY and RABBITMQ_EXCHANGE are required")
	}

	r.logger.Info("starting consumer")
	c, err := utils.StartConsumer(r.logger, "stream-automator")
	if err != nil {
		r.logger.Fatal("Error starting consumer", zap.Error(err))
	}
	defer func(consumer *utils.Consumer) {
		err := consumer.Shutdown()
		if err != nil {
			r.logger.Fatal("Error shutting down consumer", zap.Error(err))
		}
	}(c)

	consumerHandler := adapterConsumer.NewRabbitTaskQueueConsumer(c.Channel, r.taskController, r.logger, r.ctx, queueName, consumerName, bindingKey, exchange)
	consumerUseCase := consumer.NewTaskQueueConsumer(consumerHandler)

	return consumerUseCase.StartConsumer()
}
