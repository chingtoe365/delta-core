package signalutil

import (
	"context"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/notificationutil"
	"delta-core/repository"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChangeSignalConfig struct {
	Duration   int     `mapstructure:"duration"` // in seconds
	Percentage float32 `mapstructure:"percentage"`
	IsUp       bool    `mapstructure:"isUp"`
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
	log.Print(config)
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

	// log.Printf(">> Task ID %s\n", task.ID.Hex())
	ok := cs.TaskMap.TryFetch(signalId.Hex())
	if ok {
		log.Printf(">> Already been setup\n")
		return
	}
	// update task map to include task ID
	cs.TaskMap.Update(signalId.Hex(), false)
	log.Printf(">> Signal setup with ID %s \n", signalId.Hex())

	for {
		ok := cs.TaskMap.TryFetch(signalId.Hex())
		if !ok {
			// in the case when taks id is removed from map
			// which means unsubsribe
			// it will exit
			log.Printf(">> Signal with ID %s cancelled \n", signalId.Hex())
			// client.Disconnect(250)
			return
		}
		// signal calculation logic goes here
		end := time.Now()
		duration := -time.Duration(cs.Config.Duration) * time.Second
		start := end.Add(duration)
		tss := cs.Repository.FetchRawSeries(*cs.Context, cs.Key, start, end)
		if len(tss) > 0 {
			first := tss[0]
			last := tss[len(tss)-1]
			elapsed := float64(last.Timestamp-first.Timestamp) / 1000
			log.Println(elapsed)
			coverage := float64(elapsed) / float64(cs.Config.Duration)
			if coverage > 0.7 {
				if cs.Config.IsUp {
					if (last.Value-first.Value)/first.Value-float64(cs.Config.Percentage) > 0 {
						slog.Info("Up signal detected!!!!!!")
						// up signal detected, do something
						cs.Notify(fmt.Sprintf("Go up more than %v%% within %v seconds", cs.Config.Percentage*100, cs.Config.Duration))
						// frozen for another period of time of same length to duration
						time.Sleep(time.Duration(cs.Config.Duration) * time.Second)
					}
				}
				if !cs.Config.IsUp {
					if (last.Value-first.Value)/first.Value+float64(cs.Config.Percentage) < 0 {
						slog.Info("Down signal detected!!!!!!")
						// down signal detected, do something
						cs.Notify(fmt.Sprintf("Go down more than %v%% within %v seconds", cs.Config.Percentage*100, cs.Config.Duration))
						// frozen for another period of time of same length to duration
						time.Sleep(time.Duration(cs.Config.Duration) * time.Second)
					}
				}
			}
			// log.Print("ding.........")
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
