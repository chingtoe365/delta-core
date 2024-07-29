package signalutil

import (
	"context"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/notificationutil"
	"delta-core/repository"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChangeSignalConfig struct {
	duration   int
	percentage float32
	isUp       bool
}

type ChangeSignaler struct {
	SignalId   primitive.ObjectID
	Key        string
	Config     ChangeSignalConfig
	TaskMap    domain.SafeTaskMap
	Repository *repository.MarketRepository
	Context    *context.Context
	Env        *bootstrap.Env
	Profile    *domain.Profile
}

func NewChangeSignaler(signalId primitive.ObjectID, key string, config ChangeSignalConfig, taskMap domain.SafeTaskMap, repository *repository.MarketRepository, ctx context.Context, env *bootstrap.Env, profile *domain.Profile) *ChangeSignaler {
	return &ChangeSignaler{
		SignalId:   signalId,
		Key:        key,
		Config:     config,
		TaskMap:    taskMap,
		Repository: repository,
		Context:    &ctx,
		Env:        env,
		Profile:    profile,
	}
}

func (cs *ChangeSignaler) Roll(signalId primitive.ObjectID) {
	// create and connect clients
	// var client = mqttutil.NewMqttClient(env, profile)

	// fmt.Printf(">> Task ID %s\n", task.ID.Hex())
	ok := cs.TaskMap.TryFetch(signalId.Hex())
	if ok {
		fmt.Printf(">> Already been setup\n")
		return
	}
	// update task map to include task ID
	cs.TaskMap.Update(signalId.Hex(), false)

	for {
		ok := cs.TaskMap.TryFetch(signalId.Hex())
		if !ok {
			// in the case when taks id is removed from map
			// which means unsubsribe
			// it will exit
			fmt.Printf(">> Signal with ID %s cancelled \n", signalId.Hex())
			// client.Disconnect(250)
			return
		}
		// signal calculation logic goes here
		// var endTs = time.Now().Unix() - int(cs.Config["duration"])
		// startTs := time.Now().Unix() - int64(cs.Config.duration)
		end := time.Now()
		start := end.Add(time.Duration(cs.Config.duration * 1000 * 1000 * -1))
		tss := cs.Repository.FetchRawSeries(*cs.Context, cs.Key, start, end)
		// tsPointtsPoints[len(tsPoints)-1]
		first := tss[0]
		last := tss[len(tss)-1]
		// if last.Time
		elapsed := (last.Timestamp - first.Timestamp) / 1000 / 1000
		coverage := float64(elapsed) / float64(cs.Config.duration)
		if coverage > 0.95 {
			if cs.Config.isUp {
				if (last.Value-first.Value)/first.Value-float64(cs.Config.percentage) > 0 {
					// up signal detected, do something
					cs.Notify(fmt.Sprintf("Go up more than %v%% within %v seconds", cs.Config.percentage*100, cs.Config.duration))
				}
			}
			if !cs.Config.isUp {
				if (last.Value-first.Value)/first.Value+float64(cs.Config.percentage) < 0 {
					// down signal detected, do something
					cs.Notify(fmt.Sprintf("Go down more than %v%% within %v seconds", cs.Config.percentage*100, cs.Config.duration))
				}
			}

		}
		// otherwise keep listening
		time.Sleep(50 * time.Millisecond)
	}
}

func (cs *ChangeSignaler) Notify(shortDesc string) {
	var a = domain.Alert{
		TradeItem: cs.Key,
		Signal: domain.Signal{
			Short:       shortDesc,
			Description: shortDesc,
		},
		Time: time.Now(),
	}
	// a.ParseIn(shortDesc, msg.Topic())
	notificationutil.SendMail(cs.Env, cs.Profile.Email, a.FormatEmail())
}

func (cs *ChangeSignaler) Remove() {

}
