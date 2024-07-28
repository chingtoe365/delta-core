package usecase

import (
	"context"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/signalutil"
	"delta-core/repository"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignalSetupUsecase struct {
	// TaskMap chan map[string]bool
	TaskMap domain.SafeTaskMap
}

func NewSignalSetupUsecase() *SignalSetupUsecase {
	return &SignalSetupUsecase{
		TaskMap: domain.SafeTaskMap{
			Status: make(map[string]bool), Unlocked: make(chan int),
		},
	}
}

func (ssu *SignalSetupUsecase) MakeMarketSignaler(signalId primitive.ObjectID, signalerKey string, signalerType string, singalerConfig map[string]string, repository *repository.MarketRepository, context context.Context, env *bootstrap.Env, profile *domain.Profile) (domain.IMarketSignaler, error) {
	var signaler domain.IMarketSignaler
	switch signalerType {
	case "change":
		var signalChangeConfig signalutil.ChangeSignalConfig
		signalerJsonStr, err := json.Marshal(singalerConfig)
		if err != nil {
			return &signalutil.ChangeSignaler{}, err
		}
		err = json.Unmarshal([]byte(signalerJsonStr), &signalChangeConfig)
		if err != nil {
			return &signalutil.ChangeSignaler{}, err
		}
		signaler = signalutil.NewChangeSignaler(signalId, signalerKey, signalChangeConfig, ssu.TaskMap, repository, context, env, profile)
	default:
		signaler = &signalutil.ChangeSignaler{}
	}
	return signaler, nil
}

func (ssu *SignalSetupUsecase) RemoveMarketSignaler(signalId string) {
	// remove task id in task map to kill goroutine
	ok := ssu.TaskMap.TryFetch(signalId)
	if !ok {
		fmt.Printf("Signal deleted \n")
		return
	}
	ssu.TaskMap.Update(signalId, true)
}

// func setupSignals(config domain.MarketSignalSetupRequest) {
// 	if config.Cateogry == "change" {

// 	}
// }
