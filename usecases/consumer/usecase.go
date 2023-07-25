package consumer

type TaskQueueConsumer struct {
	adapter TaskQueueConsumerAdapter
}

func NewTaskQueueConsumer(adapter TaskQueueConsumerAdapter) *TaskQueueConsumer {
	return &TaskQueueConsumer{adapter: adapter}
}

func (tqc *TaskQueueConsumer) StartConsumer() []error {
	return tqc.adapter.ConsumeTasks()
}
