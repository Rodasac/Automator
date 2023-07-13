package consumer

import (
	"automator-go/adapters/controllers/tasks"
	"automator-go/entities/models"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitTaskQueueConsumer struct {
	taskController *tasks.TaskController
	ch             *amqp.Channel
	chName         string
	consumerName   string
	ctx            context.Context
}

func NewRabbitTaskQueueConsumer(
	taskController *tasks.TaskController,
	ch *amqp.Channel,
	chName string,
	consumerName string,
	ctx context.Context,
) RabbitTaskQueueConsumer {
	return RabbitTaskQueueConsumer{
		taskController: taskController,
		ch:             ch,
		chName:         chName,
		consumerName:   consumerName,
		ctx:            ctx,
	}
}

func (t RabbitTaskQueueConsumer) ConsumeTasks() error {
	q, err := t.ch.QueueDeclare(
		t.chName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error declaring queue: %w", err)
	}

	msgs, err := t.ch.Consume(
		q.Name,
		t.consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error consuming queue: %w", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received task: %s", d.Body)
			var taskToProcess models.Task
			err = json.Unmarshal(d.Body, &taskToProcess)
			if err != nil {
				log.Printf("Error unmarshalling task: %s", err)
				err := d.Nack(false, false)
				if err != nil {
					log.Printf("Error nacknowledging message: %s", err)
				}

				continue
			}

			err = t.taskController.ProcessTask(&taskToProcess)
			if err != nil {
				log.Printf("Error processing task: %s", err)
				err := d.Nack(false, false)
				if err != nil {
					log.Printf("Error nacknowledging message: %s", err)
				}

				continue
			}

			err = d.Ack(true)
			if err != nil {
				log.Printf("Error acknowledging message: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for tasks. To exit press CTRL+C")
	<-t.ctx.Done()

	return nil
}
