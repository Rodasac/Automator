package consumer

import (
	"automator-go/robot/adapters/controllers/tasks"
	"automator-go/robot/entities/models"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type RabbitTaskQueueConsumer struct {
	ch             *amqp.Channel
	taskController *tasks.TaskController
	logger         *otelzap.LoggerWithCtx
	ctx            context.Context
	queueName      string
	consumerName   string
}

func NewRabbitTaskQueueConsumer(
	ch *amqp.Channel,
	taskController *tasks.TaskController,
	logger *otelzap.LoggerWithCtx,
	ctx context.Context,
	queueName string,
	consumerName string,
) RabbitTaskQueueConsumer {
	return RabbitTaskQueueConsumer{
		ch:             ch,
		taskController: taskController,
		logger:         logger,
		ctx:            ctx,
		queueName:      queueName,
		consumerName:   consumerName,
	}
}

func (t RabbitTaskQueueConsumer) startConsumer() (<-chan amqp.Delivery, error) {
	t.logger.Debug("Queue bound to Exchange, starting Consume", zap.String("consumerTag", t.consumerName))
	deliveries, err := t.ch.Consume(
		t.queueName,
		t.consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}

	return deliveries, nil
}

func (t RabbitTaskQueueConsumer) ConsumeTasks() []error {
	deliveries, err := t.startConsumer()
	if err != nil {
		return []error{err}
	}

	go func() {
		for d := range deliveries {
			delivery := d
			t.logger.Debug("received message", zap.ByteString("body", delivery.Body))
			var taskToProcess models.Task
			err = json.Unmarshal(delivery.Body, &taskToProcess)
			if err != nil {
				t.logger.Error("Error unmarshalling task", zap.Error(err))
				err := delivery.Nack(false, false)
				if err != nil {
					t.logger.Error("Error nacknowledging message", zap.Error(err))
				}

				continue
			}

			// Process tasks concurrently
			go func() {
				err = t.taskController.ProcessTask(&taskToProcess)
				if err != nil {
					t.logger.Error("Error processing task", zap.Error(err))
					err := delivery.Nack(false, false)
					if err != nil {
						t.logger.Error("Error nacknowledging message", zap.Error(err))
					}

					return
				}

				err = delivery.Ack(false)
				if err != nil {
					t.logger.Error("Error acknowledging message", zap.Error(err))
				}
			}()
		}
	}()

	t.logger.Info("[*] Waiting for tasks. To exit press CTRL+C")
	<-t.ctx.Done()

	return nil
}
