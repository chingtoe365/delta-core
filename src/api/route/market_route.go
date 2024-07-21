package route

import (
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/mongo"
	"delta-core/repository"

	"github.com/gin-gonic/gin"
)

func NewMarketRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {

	mc := &controller.MarketController{
		MarketRepository: repository.NewRedisMarketRepository(),
	}

	group.GET("/get-all-trade-items", mc.GetAllTradeItems)
	group.GET("/get-all-trade-signals", mc.GetAllTradeSignals)
	group.GET("/get-series", mc.GetSeries)
	// group.GET("/get-quote", mc.GetQuote)

}
