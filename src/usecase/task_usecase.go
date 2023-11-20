package usecase

import (
	"context"
	"time"

	"delta-core/domain"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (tu *taskUsecase) Create(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Create(ctx, task)
}

func (tu *taskUsecase) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchByUserID(ctx, userID)
}
func (tu *taskUsecase) FetchAll(c context.Context) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchAll(ctx)
}
func (tu *taskUsecase) Delete(c context.Context, task *domain.Task) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Delete(ctx, task)
}

func (tu *taskUsecase) FetchByTaskID(c context.Context, taskId string) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchById(ctx, taskId)
}
