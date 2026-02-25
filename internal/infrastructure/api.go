// Package infrastructure provides the implementation of infrastructure components for the application.
// It includes functions to create and configure the API engine, database connections, and other necessary services.
// This package serves as the bridge between the application's core logic and external systems.
// It ensures that all dependencies are properly initialized and injected into the application.
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/logger"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/user-service/internal/api"
)

// CreateAPIConfig initializes and returns the API configuration.
// Returns:
//   - *api.Config: A pointer to the initialized API configuration.
func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	common.HandlerError(err)

	return cfg
}

// CreatAPI initializes and returns the API engine.
// It sets up all necessary dependencies including database connections, JWT providers, and other services.
// Returns:
//   - api.Engine: The initialized API engine.
func CreateAPI() api.Engine {
	logger.SetLogLevel()

	// Initialize API configuration
	cfg := CreateAPIConfig()

	// initialize redis client
	redisClient := CreateRedisCon()

	// initialize sql db client
	dbClient := CreateSQLDBAndMigration()

	// initialize other dependencies
	jwtGenerator, jwtValidator := CreateJWTProviders()
	app := gin.New()

	apiEngine := api.New(&api.EngineOpts{
		Engine:      app,
		Cfg:         cfg,
		RedisClient: redisClient,
		SqlDB:       dbClient,
		// Initialize other dependencies
		RandomCodeGen:   utils.NewCodeGenerator(),
		PasswordHashing: utils.NewPasswordHashing(),
		JWTGenerator:    jwtGenerator,
		JWTValidator:    jwtValidator,
	})

	return apiEngine
}
