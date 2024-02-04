package route

import (
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/mongo"

	"github.com/gin-gonic/gin"
)

func NewMarketRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {

	mc := &controller.MarketController{}

	group.GET("/get-all-trade-items", mc.GetAllTradeItems)
	group.GET("/get-all-trade-signals", mc.GetAllTradeSignals)
	group.GET("/get-quote", mc.GetQuote)

}
