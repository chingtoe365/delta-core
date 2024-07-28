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

func NewMarketRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pr := repository.NewUserRepository(db, domain.CollectionUser)
	mc := &controller.MarketController{
		MarketRepository:   repository.NewMarketRepository(db, domain.CollectionMarketSignal),
		SignalSetupUsecase: usecase.NewSignalSetupUsecase(),
		ProfileUsecase:     usecase.NewProfileUsecase(pr, timeout),
		Env:                env,
	}

	group.GET("/get-all-trade-items", mc.GetAllTradeItems)
	group.GET("/get-all-trade-signals", mc.GetAllTradeSignals)
	group.GET("/get-all-trade-signal-categories", mc.GetAllTradeSignalCategories)
	group.GET("/get-series", mc.GetSeries)
	// group.POST("/setup-trade-signal", mc.SetupTradeSignals)
	// group.GET("/get-quote", mc.GetQuote)

}
