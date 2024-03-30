package controller

import (
	_ "example.com/sample/microsvc/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type SwaggerController struct {
}

// New constructs a new Swagger Controller that allows us to add swagger and have the propper swagger annotations
func New() SwaggerController {
	return SwaggerController{}
}

// AttachRoutes implements router.Controller for SwaggerController. It attaches this controller's routes to the microservice's router.
func (t SwaggerController) AttachRoutes(rtr *echo.Echo) {
	rtr.GET("/api/swagger/*", echoSwagger.WrapHandler)
}
