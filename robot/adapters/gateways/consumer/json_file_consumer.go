package consumer

import (
	"automator-go/robot/adapters/controllers/tasks"
	"automator-go/robot/entities/models"
	"encoding/json"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"os"
	"sync"
)

type TaskQueueConsumerFromJSONFile struct {
	taskController *tasks.TaskController
	logger         *otelzap.LoggerWithCtx
}

func NewTaskQueueConsumerFromJSONFile(
	taskController *tasks.TaskController,
	logger *otelzap.LoggerWithCtx,
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

	wg := sync.WaitGroup{}
	wg.Add(len(tasksToProcess))
	errorsChan := make(chan error, len(tasksToProcess))
	for _, task := range tasksToProcess {
		taskToProcess := task
		go func() {
			defer func() {
				wg.Done()
				t.logger.Debug("Finished processing task", zap.String("task_id", taskToProcess.Id))
			}()

			t.logger.Info("Processing task", zap.String("task_id", taskToProcess.Id))
			err := t.taskController.ProcessTask(&taskToProcess)
			if err != nil {
				errorsChan <- fmt.Errorf("error processing task: %s: %w", taskToProcess.Id, err)
			}
		}()
	}
	wg.Wait()
	close(errorsChan)

	t.logger.Debug("Finished processing all tasks")

	errors := make([]error, 0)
	for {
		err, ok := <-errorsChan
		if !ok {
			break
		}

		errors = append(errors, err)
	}

	return errors
}
