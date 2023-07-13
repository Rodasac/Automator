package consumer

type TaskQueueConsumerAdapter interface {
	ConsumeTasks() error
}

type TaskQueueConsumerUseCase interface {
	StartConsumer() error
}
