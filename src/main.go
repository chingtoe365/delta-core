package main

import (
	"time"

	route "delta-core/api/route"
	"delta-core/bootstrap"
	docs "delta-core/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title DeltaTrade Core Swagger API
// @version 1.0
// @description Swagger API documentation for DeltaTrade core service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				"Type 'Bearer TOKEN' to correctly set the API Key"
func main() {

	app := bootstrap.App()

	env := app.Env

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

	route.Setup(env, timeout, db, ginApp)

	docs.SwaggerInfo.BasePath = "/api/v1"

	ginApp.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	ginApp.Run(env.ServerAddress)

}
