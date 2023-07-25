package consumer

import (
	"automator-go/adapters/controllers/tasks"
	"automator-go/entities/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type TaskQueueConsumerFromJSONFile struct {
	taskController *tasks.TaskController
	logger         *zap.Logger
}

func NewTaskQueueConsumerFromJSONFile(
	taskController *tasks.TaskController,
	logger *zap.Logger,
) TaskQueueConsumerFromJSONFile {
	return TaskQueueConsumerFromJSONFile{taskController: taskController, logger: logger}
}

func (t TaskQueueConsumerFromJSONFile) ConsumeTasks() []error {
	t.logger.Debug("Opening file")
	file, err := os.ReadFile("./tasks_test.json")
	if err != nil {
		return []error{err}
	}

	t.logger.Debug("Unmarshalling file")
	var tasksToProcess []models.Task
	err = json.Unmarshal(file, &tasksToProcess)
	if err != nil {
		return []error{err}
	}

	errors := make([]error, 0)
	for _, task := range tasksToProcess {
		t.logger.Info("Processing task", zap.String("task_id", task.Id))
		err = t.taskController.ProcessTask(&task)
		if err != nil {
			errors = append(errors, fmt.Errorf("error processing task: %s: %w", task.Id, err))
		}
	}

	return errors
}
