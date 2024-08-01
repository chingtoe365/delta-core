package route

import (
	"time"

	"delta-core/api/middleware"
	"delta-core/bootstrap"
	"delta-core/mongo"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	generalRouter := gin.Group(env.BasePath)
	{
		publicRouter := generalRouter.Group("")
		// All Public APIs
		NewSignupRouter(env, timeout, db, publicRouter)
		NewLoginRouter(env, timeout, db, publicRouter)
		NewRefreshTokenRouter(env, timeout, db, publicRouter)
		NewTestRouter(env, timeout, db, publicRouter)
		NewMarketRouter(env, timeout, db, publicRouter)
		NewNewsRouter(env, timeout, db, publicRouter)

		protectedRouter := generalRouter.Group("")
		// Middleware to verify AccessToken
		protectedRouter.Use(middleware.JwtAuthMiddleware(env))
		// All Private APIs
		NewProfileRouter(env, timeout, db, protectedRouter)
		NewTaskRouter(env, timeout, db, protectedRouter)
		NewSignalRouter(env, timeout, db, protectedRouter)
	}
}
