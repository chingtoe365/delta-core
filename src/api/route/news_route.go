package route

import (
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/mongo"

	"github.com/gin-gonic/gin"
)

func NewNewsRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {

	nc := &controller.NewsController{}

	group.GET("/get-news", nc.Get)
}
