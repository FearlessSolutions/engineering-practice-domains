package main

import (
	"log"

	"example.com/sample/commonlib/database"
	"example.com/sample/commonlib/logger"
	"example.com/sample/commonlib/router"
	"example.com/sample/commonlib/router/middleware"
	loglevelcontroller "example.com/sample/commonlib/sharedfeatures/loglevel/controller"
	sampleadapter "example.com/sample/microsvc/features/sample/adapter"
	samplecontroller "example.com/sample/microsvc/features/sample/controller"
	swaggercontroller "example.com/sample/microsvc/features/swagger/controller"
	"example.com/sample/microsvc/options"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// PrepareSubsystems prepares the set of global systems that other parts of the microservice depend on,
// namely the global logger (logger.Log) and the global configuration registry (options.Registry)
func PrepareSubsystems() *sqlx.DB {
	// Attempt to load .env file, errors are inconsequential
	_ = godotenv.Load()

	// Set up config registry
	configSetupErr := options.InitRegistry()
	if configSetupErr != nil {
		log.Fatalf("Could not set up configuration registry: %v", configSetupErr)
	}

	// Set up logger
	loggerSetupErr := logger.InitLoggerFromConfig(*options.Registry)
	if loggerSetupErr != nil {
		log.Fatal("Could not set up logger!", loggerSetupErr)
	}

	// Set up database connection
	db, dbConnectErr := database.ConnectFromConfig(*options.Registry)
	if dbConnectErr != nil {
		logger.Log.Fatal("Database connection irreparably failed!", zap.Error(dbConnectErr))
	}

	// Verify the connection is established
	database.MustBeConnected(db)

	return db
}

// Bootstrap constructs the microservice's controllers and middleware, then creates a router and attaches
// the controllers and middleware to it
func Bootstrap(db *sqlx.DB) router.Router {
	controllers := CreateControllers()
	appMiddleware := CreateMiddleware(db)

	appRouter := router.New()
	appRouter.AttachMiddleware(appMiddleware)
	appRouter.AttachControllers(controllers)

	return appRouter
}

// CreateControllers constructs all the rest controllers in the microservice
func CreateControllers() []router.Controller {
	return []router.Controller{
		sample(),
		logLevelAdjust(),
		swagger(),
	}
}

// CreateMiddleware constructs all the middleware the microservice will use
func CreateMiddleware(db *sqlx.DB) []echo.MiddlewareFunc {
	var appMiddleware []echo.MiddlewareFunc
	appMiddleware = append(appMiddleware, middleware.StandardMiddleware(*options.Registry, db)...)

	return appMiddleware
}

// sample constructs the sample loglevelcontroller (controller.SampleController)
func sample() samplecontroller.SampleController {
	greetingReader := sampleadapter.DatabaseGreetingReader{}
	greetingWriter := sampleadapter.DatabaseGreetingWriter{}
	return samplecontroller.New(greetingReader, greetingWriter)
}

// swagger creates a swagger controller used for generating the swagger for the project
func swagger() swaggercontroller.SwaggerController {
	return swaggercontroller.New()
}

// logLevelAdjust constructs the shared log level adjustment controller (controller.LogLevelController)
func logLevelAdjust() loglevelcontroller.LogLevelController {
	return loglevelcontroller.New()
}
