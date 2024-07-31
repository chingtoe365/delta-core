package route

import (
	"context"
	"log"
	"log/slog"
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/mongo"
	"delta-core/repository"
	"delta-core/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func NewSignalRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pr := repository.NewUserRepository(db, domain.CollectionUser)
	mc := &controller.MarketController{
		MarketRepository:   repository.NewMarketRepository(db, domain.CollectionMarketSignal),
		SignalSetupUsecase: usecase.NewSignalSetupUsecase(),
		ProfileUsecase:     usecase.NewProfileUsecase(pr, timeout),
		Env:                env,
	}
	var signalCollection = db.Collection(domain.CollectionMarketSignal)
	var signals []domain.MarketSignalDto

	cursor, err := signalCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Cannot fetch all tasks, error %v", err)
	}

	cursor.All(context.TODO(), &signals)

	slog.Info("Initializing signals setup")
	mc.SignalSetupUsecase.InitialiseSignalsSetup(context.TODO(), env, mc.ProfileUsecase, mc.MarketRepository, signals)

	group.GET("/list-signals", mc.ListSubscribedTradeSignals)
	group.POST("/setup-signal", mc.SetupTradeSignals)
	group.DELETE("/delete-signal", mc.DeleteTradeSignals)
	// group.GET("/get-quote", mc.GetQuote)

}
