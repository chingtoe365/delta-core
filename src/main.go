package main

import (
	"log/slog"
	"os"
	"time"

	route "delta-core/api/route"
	"delta-core/bootstrap"
	docs "delta-core/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func ApiEnvMiddleware(env *bootstrap.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("env", env)
		c.Next()
	}
}

// @title DeltaTrade Core Swagger API
// @version 1.0
// @description Swagger API documentation for DeltaTrade core service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				"Type 'Bearer TOKEN' to correctly set the API Key"
func main() {

	app := bootstrap.App()

	env := app.Env

	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// logger := slog.NewLogLogger(slog.NewTextHandler())
	// slog.SetDefault(logger)
	// defaultLogger := log.Default()
	// defaultLogger.SetOutput(os.Stdout)
	// slog.Info("Hello from the std-log package!")
	// logger := log.New(
	// 	os.Stderr,
	// 	"",
	// 	// "Go application: ",
	// 	log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile,
	// )
	// slog.SetDefault(logger)
	if env.LogVerbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})))
	}

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	if env.AppEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else if env.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {

	}

	ginApp := gin.Default()

	// attaching global middle to pass env object into context
	ginApp.Use(ApiEnvMiddleware(env))

	route.Setup(env, timeout, db, ginApp)

	docs.SwaggerInfo.BasePath = "/api/v1"

	ginApp.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	ginApp.Run(env.ServerAddress)

}
