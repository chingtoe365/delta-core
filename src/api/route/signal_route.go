package route

import (
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/mongo"
	"delta-core/repository"
	"delta-core/usecase"

	"github.com/gin-gonic/gin"
)

func NewSignalRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pr := repository.NewUserRepository(db, domain.CollectionUser)
	mc := &controller.MarketController{
		MarketRepository:   repository.NewMarketRepository(db, domain.CollectionMarketSignal),
		SignalSetupUsecase: usecase.NewSignalSetupUsecase(),
		ProfileUsecase:     usecase.NewProfileUsecase(pr, timeout),
		Env:                env,
	}

	group.GET("/list-signals", mc.ListSubscribedTradeSignals)
	group.POST("/setup-signal", mc.SetupTradeSignals)
	group.DELETE("/delete-signal", mc.DeleteTradeSignals)
	// group.GET("/get-quote", mc.GetQuote)

}
