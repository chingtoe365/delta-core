package route

import (
	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/mongo"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tc := &controller.TestController{
		Env: env,
	}
	group.GET("/test", tc.Test)
}
