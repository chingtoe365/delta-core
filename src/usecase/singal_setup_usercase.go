package usecase

import (
	"context"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/internal/signalutil"
	"delta-core/repository"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
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

func (ssu *SignalSetupUsecase) MakeMarketSignaler(signalId primitive.ObjectID, signalerKey string, signalerType string, singalerConfig map[string]interface{}, repository *repository.MarketRepository, context context.Context, env *bootstrap.Env, profile *domain.Profile) (domain.IMarketSignaler, error) {
	var signaler domain.IMarketSignaler
	fmt.Println(singalerConfig)
	switch signalerType {
	case "change":
		var signalChangeConfig signalutil.ChangeSignalConfig
		config := &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &signalChangeConfig,
		}

		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			panic(err)
		}
		log.Println("Before decoding")
		err = decoder.Decode(singalerConfig)
		if err != nil {
			println(err)
			log.Println("error happen when mapping")
			log.Fatal("Cannot convert map to struct")
		}
		log.Println(signalChangeConfig)

		signaler = signalutil.NewChangeSignaler(signalId, signalerKey, signalChangeConfig, ssu.TaskMap, repository, context, env, profile)
	default:
		signaler = &signalutil.ChangeSignaler{}
	}
	fmt.Println(signaler)
	return signaler, nil
}

func (ssu *SignalSetupUsecase) RemoveMarketSignaler(signalId string) {
	// remove task id in task map to kill goroutine
	ok := ssu.TaskMap.TryFetch(signalId)
	log.Print(ok)
	if !ok {
		fmt.Printf("Signal deleted \n")
		return
	}
	ssu.TaskMap.Update(signalId, true)
}

func (ssu *SignalSetupUsecase) InitialiseSignalsSetup(ctx context.Context, env *bootstrap.Env, puc domain.ProfileUsecase, repository *repository.MarketRepository, signals []domain.MarketSignalDto) {
	// unlock task map channel first
	go ssu.TaskMap.Unlock()
	for _, item := range signals {
		fmt.Printf("Starting signal %s ID: %s\n", item.Id, item.Id.Hex())
		// fmt.Println(item.UserID.Hex())
		profile, err := puc.GetProfileByID(ctx, item.UserId.Hex())
		if err != nil {
			panic(err)
		}
		fmt.Println("Before making signals")
		signaler, err := ssu.MakeMarketSignaler(
			item.Id, item.SignalMeta.Key, string(item.SignalMeta.Type), item.SignalMeta.Config, repository, ctx, env, profile)
		if err != nil {
			// ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
			log.Fatal(err)
			// return
			log.Printf("Cannot initialize  signal %s ID: %s\n", item.Id, item.Id.Hex())
		}
		go signaler.Roll(item.Id)
		// go ssu.Subscribe(env, item, profile)
	}
}

// func setupSignals(config domain.MarketSignalSetupRequest) {
// 	if config.Cateogry == "change" {

// 	}
// }
