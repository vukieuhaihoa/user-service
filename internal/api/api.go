// Package api provides the HTTP API server implementation for the bookmark management application.
// It handles routing, middleware, and server lifecycle management using the Gin web framework.
package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vukieuhaihoa/user-service/docs"
	"gorm.io/gorm"

	middleware "github.com/vukieuhaihoa/bookmark-libs/middlewares"

	"github.com/vukieuhaihoa/bookmark-libs/ratelimit"

	healthCheckHandler "github.com/vukieuhaihoa/user-service/internal/app/handler/healthcheck"
	healthCheckRepository "github.com/vukieuhaihoa/user-service/internal/app/repository/healthcheck"
	healthCheckService "github.com/vukieuhaihoa/user-service/internal/app/service/healthcheck"

	userHandler "github.com/vukieuhaihoa/user-service/internal/app/handler/user"
	userRepository "github.com/vukieuhaihoa/user-service/internal/app/repository/user"
	userService "github.com/vukieuhaihoa/user-service/internal/app/service/user"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/jwtutils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/utils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/validators"
)

var registerValidationsOnce sync.Once

// Engine defines the contract for the HTTP server engine.
// It provides methods to start and manage the API server lifecycle.
type Engine interface {
	// Start initializes and starts the HTTP server.
	// Returns an error if the server fails to start or encounters a runtime error.
	Start() error

	// ServeHTTP allows the engine to handle HTTP requests.
	// It satisfies the http.Handler interface.
	// Parameters:
	//   - w: The http.ResponseWriter to write the response
	//   - r: The incoming http.Request to be handled
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// api is the concrete implementation of the Engine interface.
// It encapsulates the Gin engine and provides HTTP server functionality.
type api struct {
	// app is the underlying Gin engine instance that handles HTTP routing and middleware
	app *gin.Engine

	// cfg holds the configuration settings for the API server
	cfg *Config

	// redisClient is the Redis client used for caching and session management
	redisClient *redis.Client

	// randomCodeGen is the code generator used for generating random codes
	randomCodeGen utils.CodeGenerator

	passwordHashing utils.PasswordHashing

	db *gorm.DB

	jwtGenerator jwtutils.JWTGenerator

	jwtValidator jwtutils.JWTValidator
}

type EngineOpts struct {
	Engine          *gin.Engine
	Cfg             *Config
	RedisClient     *redis.Client
	SqlDB           *gorm.DB
	RandomCodeGen   utils.CodeGenerator
	PasswordHashing utils.PasswordHashing
	JWTGenerator    jwtutils.JWTGenerator
	JWTValidator    jwtutils.JWTValidator
}

// New creates a new instance of the API engine with the provided options.
// It initializes the Gin engine, configures routes, and sets up middleware.
//
// Parameters:
//   - opts: EngineOpts containing dependencies and configurations for the API engine
//
// Returns:
//   - Engine: An instance of the API engine implementing the Engine interface
func New(opts *EngineOpts) Engine {
	a := &api{
		app:             opts.Engine,
		cfg:             opts.Cfg,
		redisClient:     opts.RedisClient,
		randomCodeGen:   opts.RandomCodeGen,
		passwordHashing: opts.PasswordHashing,
		db:              opts.SqlDB,
		jwtGenerator:    opts.JWTGenerator,
		jwtValidator:    opts.JWTValidator,
	}

	a.registerValidations()
	a.registerRoutes()

	return a
}

// Start initializes and starts the HTTP server.
// It listens on the configured application port and begins handling incoming requests.
// defaults to port 8080 if not specified in the configuration.
// Returns:
//   - error: An error if the server fails to start, nil otherwise
func (a *api) Start() error {
	return a.app.Run(a.cfg.AppPort)
}

// ServeHTTP allows the api struct to satisfy the http.Handler interface.
// It delegates HTTP requests to the underlying Gin engine.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response
//   - r: The incoming http.Request to be handled
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// registerRoutes configures all HTTP routes and their handlers for the API.
// This method initializes services and handlers, then registers them with the Gin engine.
// Currently registers the password generation endpoint.
func (a *api) registerRoutes() {
	allHandler := a.registerHandlers()
	allMiddlewares := a.registerMiddlewares()

	// Swagger info setup
	docs.SwaggerInfo.Host = a.cfg.AppHostName

	// Register health check endpoint
	a.app.GET("/health-check", allHandler.healthCheckHandler.Check)
	// Register Swagger documentation endpoint
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API version group
	v1 := a.app.Group("/v1")
	v1.Use(allMiddlewares.rateLimitMiddleware.RateLimit(middleware.RateLimitIPKey)) // Apply rate limiting middleware to all /v1 routes
	{
		v1.POST("/users/register", allHandler.userHandler.CreateUser)

		v1.POST("/users/login", allHandler.userHandler.Login)

	}

	v1Private := a.app.Group("/v1")
	v1Private.Use(allMiddlewares.jwtAuth.JWTAuth())
	v1Private.Use(allMiddlewares.rateLimitMiddleware.RateLimit(middleware.RateLimitUserIDKey)) // Apply rate limiting middleware to all /v1 routes for authenticated users
	{
		v1Private.GET("/self/info", allHandler.userHandler.GetProfile)
		v1Private.PUT("/self/info", allHandler.userHandler.UpdateProfile)
	}
}

func (a *api) registerValidations() {
	registerValidationsOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			// Register custom validation functions here
			v.RegisterValidation("password_strength", validators.PasswordStrength)
		}
	})
}

// handlers aggregates all HTTP handlers for different API endpoints.
type handlers struct {
	healthCheckHandler healthCheckHandler.Handler
	userHandler        userHandler.Handler
}

// registerHandlers initializes and returns all handler instances used in the API.
func (a *api) registerHandlers() *handlers {
	healthCheckRepo := healthCheckRepository.NewHealthCheckRepository(a.redisClient, a.db)
	healthCheckSvc := healthCheckService.NewHealthCheckService(a.cfg.ServiceName, a.cfg.InstanceID, healthCheckRepo)
	healthCheckHandler := healthCheckHandler.NewHealthCheckHandler(healthCheckSvc)

	userRepo := userRepository.NewUserRepository(a.db)
	userSvc := userService.NewUserService(userRepo, a.passwordHashing, a.jwtGenerator)
	userHandler := userHandler.NewUserHandler(userSvc)

	return &handlers{
		healthCheckHandler: healthCheckHandler,
		userHandler:        userHandler,
	}
}

// middlewares aggregates all middleware instances used in the API.
type middlewares struct {
	jwtAuth             middleware.JWTAuth
	rateLimitMiddleware middleware.RateLimit
}

// registerMiddlewares configures and returns all middleware instances used in the API.
func (a *api) registerMiddlewares() *middlewares {
	jwtAuth := middleware.NewJWTAuth(a.jwtValidator)

	rateLimitRepo := ratelimit.NewRedisRepo(a.redisClient)
	rateLimitMiddleware := middleware.NewRateLimit(rateLimitRepo)

	return &middlewares{
		jwtAuth:             jwtAuth,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}
