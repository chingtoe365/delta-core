package domain

import (
	"context"
	"delta-core/bootstrap"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     primitive.ObjectID `bson:"_id" json:"-"`
	Title  string             `bson:"title" form:"title" binding:"required" json:"title"`
	UserID primitive.ObjectID `bson:"userID" json:"-"`
}

type SafeTaskMap struct {
	Status   map[string]bool
	Unlocked chan int
}

func (stm *SafeTaskMap) Unlock() {
	stm.Unlocked <- 1
}

func (stm *SafeTaskMap) Update(taskId string, remove bool) {
	for {
		select {
		case <-stm.Unlocked:
			// lock it then update the map
			if remove {
				delete(stm.Status, taskId)
			} else {
				stm.Status[taskId] = true
			}
			// log.Printf(">> Task map table updated\n")
			// unlocking through here so it get locked
			// when map is processing in this select case
			go stm.Unlock()
			return
		default:
			// still locked
			// wait for 50 ms and try again
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (stm *SafeTaskMap) TryFetch(taskId string) bool {
	for {
		select {
		case <-stm.Unlocked:
			_, ok := stm.Status[taskId]
			go stm.Unlock()
			return ok
		default:
			time.Sleep((50 * time.Millisecond))
		}
	}
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchById(c context.Context, taskId string) (Task, error)
	FetchByUserID(c context.Context, userID string) ([]Task, error)
	FetchAll(c context.Context) ([]Task, error)
	Delete(c context.Context, task *Task) error
}

type TaskUsecase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
	FetchByTaskID(c context.Context, taskId string) (Task, error)
	FetchAll(c context.Context) ([]Task, error)
	Delete(c context.Context, task *Task) error
}

type SignalSubUsecase interface {
	Subscribe(env *bootstrap.Env, task Task, profile *Profile)
	Unsubscribe(task *Task)
	InitialiseSingalSubs(c context.Context, env *bootstrap.Env, puc ProfileUsecase, tasks []Task)
}
