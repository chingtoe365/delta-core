package usecase

import (
	"context"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/mqttutil"
	"log"
	"log/slog"
	"time"
)

type SignalSubUsecase struct {
	// TaskMap chan map[string]bool
	TaskMap domain.SafeTaskMap
}

func NewSignalSubUsecase(taskRepository domain.TaskRepository, timeout time.Duration) *SignalSubUsecase {
	return &SignalSubUsecase{
		TaskMap: domain.SafeTaskMap{
			Status: make(map[string]bool), Unlocked: make(chan int),
		},
	}
}

func (ssu *SignalSubUsecase) Subscribe(env *bootstrap.Env, task domain.Task, profile *domain.Profile) {
	// create and connect clients
	var client = mqttutil.NewMqttClient(env, profile)

	// log.Printf(">> Task ID %s\n", task.ID.Hex())
	ok := ssu.TaskMap.TryFetch(task.ID.Hex())
	if ok {
		log.Printf(">> Already subscribed\n")
		return
	}
	// subscribe to the topic
	token := client.Subscribe(task.Title, 1, nil)
	token.Wait()
	log.Printf(">> Subscribed with ID %s to topic: %s\n", task.ID.Hex(), task.Title)

	// update task map to include task ID
	ssu.TaskMap.Update(task.ID.Hex(), false)

	for {
		ok := ssu.TaskMap.TryFetch(task.ID.Hex())
		if !ok {
			// in the case when taks id is removed from map
			// which means unsubsribe
			// it will exit
			log.Printf(">> Subscription with ID %s for topic %s ends\n", task.ID.Hex(), task.Title)
			client.Disconnect(250)
			return
		}
		// otherwise keep listening
		time.Sleep(50 * time.Millisecond)
	}
}

func (ssu *SignalSubUsecase) Unsubscribe(task *domain.Task) {
	// remove task id in task map to kill goroutine
	ok := ssu.TaskMap.TryFetch(task.ID.Hex())
	if !ok {
		log.Printf("Have already unsubscribed\n")
		return
	}
	ssu.TaskMap.Update(task.ID.Hex(), true)
}

// called when task router start up
func (ssu *SignalSubUsecase) InitialiseSingalSubs(ctx context.Context, env *bootstrap.Env, puc domain.ProfileUsecase, tasks []domain.Task) {
	// unlock task map channel first
	log.Print("Unlocking task map first")
	go ssu.TaskMap.Unlock()
	for _, item := range tasks {
		log.Printf("Starting task %s ID: %s\n", item.Title, item.ID.Hex())
		slog.Info(item.UserID.Hex())
		profile, err := puc.GetProfileByID(ctx, item.UserID.Hex())
		if err != nil {
			panic(err)
		}
		go ssu.Subscribe(env, item, profile)
	}
}
