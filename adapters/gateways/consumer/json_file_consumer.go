package consumer

import (
	"automator-go/adapters/controllers/tasks"
	"automator-go/entities/models"
	"encoding/json"
	"os"
)

type TaskQueueConsumerFromJSONFile struct {
	taskController *tasks.TaskController
}

func NewTaskQueueConsumerFromJSONFile(taskController *tasks.TaskController) TaskQueueConsumerFromJSONFile {
	return TaskQueueConsumerFromJSONFile{taskController: taskController}
}

func (t TaskQueueConsumerFromJSONFile) ConsumeTasks() error {
	file, err := os.ReadFile("./tasks_test.json")
	if err != nil {
		return err
	}

	var tasksToProcess []models.Task
	err = json.Unmarshal(file, &tasksToProcess)
	if err != nil {
		return err
	}

	for _, task := range tasksToProcess {
		err = t.taskController.ProcessTask(&task)
		if err != nil {
			return err
		}
	}

	return nil
}
