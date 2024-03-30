package router

import (
	"example.com/sample/commonlib/config"
	"example.com/sample/commonlib/config/sharedoptions"
	"example.com/sample/commonlib/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Router is a generalized HTTP server type which can accept and respond to incoming HTTP requests
type Router struct {
	engine *echo.Echo
}

// New constructs a new Router
func New() Router {
	return Router{
		engine: echo.New(),
	}
}

// Listen causes the router to start listening to HTTP requests
func (rtr *Router) Listen(registry *config.Registry) {
	portNumber, envVarPresent := registry.Get(sharedoptions.ListenPort)
	if !envVarPresent {
		portNumber = "8080"
	}

	logger.Log.Fatal("Failed to run server!", zap.Error(rtr.engine.Start(":"+portNumber)))
}

// AttachControllers attaches routes from the provided controllers to this Router
func (rtr *Router) AttachControllers(controllers []Controller) {
	for _, controller := range controllers {
		controller.AttachRoutes(rtr.engine)
	}
}

// AttachMiddleware globally attaches middlewares to this router which run before every request
func (rtr *Router) AttachMiddleware(middlewares []echo.MiddlewareFunc) {
	for _, ware := range middlewares {
		rtr.engine.Use(ware)
	}
}

// Controller is a type which can expose HTTP routes to a Router
type Controller interface {
	AttachRoutes(router *echo.Echo)
}
