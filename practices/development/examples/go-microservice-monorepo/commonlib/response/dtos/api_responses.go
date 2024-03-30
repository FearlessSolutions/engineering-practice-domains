package dtos

import (
	"github.com/labstack/echo/v4"
)

// APIErrorHelper helps easily create a standard HTTP error response type with
// canned messages and auto-propagation of error messages to an API user
type APIErrorHelper struct {
	// The HTTP status to respond with
	Status int
	// The error that caused the endpoint failure
	Error error
	// A human-readable explanation of why the error occurred
	Description string
}

// apiError is a standard error type for REST controllers
type apiError struct {
	Description string `json:"description" validate:"required"`
	Detail      string `json:"detail" validate:"required"`
}

// emptyBody is an empty response for use with swagger
type emptyBody struct{}

// Respond generates a standard HTTP error response and sends it over the provided echo.Context with the
// appropriate HTTP response code
func (errHelp APIErrorHelper) Respond(c echo.Context) error {
	var errDetail string
	if errHelp.Error != nil {
		errDetail = errHelp.Error.Error()
	}

	returnedErr := apiError{
		Description: errHelp.Description,
		Detail:      errDetail,
	}

	return c.JSON(errHelp.Status, returnedErr)
}
