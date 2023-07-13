package tasks

import (
	"automator-go/entities/models"
	"automator-go/usecases/task"
)

type TaskController struct {
	use_case task.ProcessorUseCase
}

func NewTaskController(use_case task.ProcessorUseCase) *TaskController {
	return &TaskController{use_case: use_case}
}

func (t *TaskController) ProcessTask(taskToProcess *models.Task) error {
	return t.use_case.Process(taskToProcess)
}
