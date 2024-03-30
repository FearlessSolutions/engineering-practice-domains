package controller

import (
	"context"
	"errors"
	"net/http"

	"example.com/sample/commonlib/database"
	"example.com/sample/commonlib/logger"
	"example.com/sample/commonlib/request"
	"example.com/sample/commonlib/response"
	"example.com/sample/commonlib/router"
	"example.com/sample/microsvc/features/sample"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// SampleController is an example router.Controller implementation with a route that produces greetings for users.
// It contains both the driving (sample.Core) and driven (sample.GreetingReader) ports to fully drive the application
// from end to end.
type SampleController struct {
	greetingReader sample.GreetingReader
	greetingWriter sample.GreetingWriter

	sampleLogic sample.Core
}

// New constructs a new SampleController, accepting the driven port implementation (sample.GreetingReader) it will pass
// through the business logic's driving port (sample.Core)
func New(greetingReader sample.GreetingReader, greetingWriter sample.GreetingWriter) SampleController {
	return SampleController{
		greetingReader: greetingReader,
		greetingWriter: greetingWriter,
		sampleLogic:    sample.CoreLogic{},
	}
}

// newWithCore constructs a SampleController with a mock business logic core for testing. It doesn't accept driven ports
// because it assumes the driving port is mocked
func newWithCore(testCore sample.Core) SampleController {
	return SampleController{
		sampleLogic: testCore,
	}
}

// AttachRoutes implements router.Controller for SampleController. It attaches this controller's routes to the microservice's router.
func (t SampleController) AttachRoutes(rtr *echo.Echo) {
	// Action endpoint - performing an action with a collection
	rtr.POST("/api/v1/sample/greetings/greet", router.AutoBindAndValidate(t.ProduceGreeting))
	// Collection endpoint - adding a new entry into a collection
	rtr.POST("/api/v1/sample/greetings", router.AutoBindAndValidate(t.AddGreeting))
}

// ProduceGreeting is the controller implementation for the greeting endpoint of SampleController. It accepts a request
// body with someone's name and responds with a greeting to that person.
//
// @Summary      greeting
// @Description  you provide a name and you are greeted
// @Tags         greeting
// @Accept       json
// @Produce      json
// @Param        name body SampleGreetingRequest true "The name to address the greeting to"
// @Success      200  {object} SampleGreetingResponse "Success"
// @Failure      400 {object} dtos.apiError "The incoming JSON payload was malformed"
// @Router       /api/v1/sample/greetings/greet [post]
func (t SampleController) ProduceGreeting(ctx echo.Context, requestedGreeting SampleGreetingRequest) error {
	greetingText, greetingErr := t.sampleLogic.GiveGreeting(request.ExtractContext(ctx), requestedGreeting.Name, t.greetingReader)
	if greetingErr != nil {
		logger.Log.Error("Failed to retrieve greeting.", zap.Error(greetingErr))
		resp := response.InternalServerError(greetingErr)
		resp.Description = "Something went wrong trying to get your greeting"
		return resp.Respond(ctx)
	}

	return ctx.JSON(http.StatusOK, SampleGreetingResponse{
		Greeting: greetingText,
	})
}

// AddGreeting is the controller implementation for the "add greeting" endpoint of SampleController. It adds a new
// greeting to the existing set of greetings.
//
// @Summary      Add a greeting
// @Description  adds a greeting
// @Tags         greeting
// @Accept       json
// @Produce      json
// @Param        greeting body NewGreetingRequest true "'greeting' is the key to the value to be added into the database"
// @Success      201 {object} dtos.emptyBody "Successfully created the greeting"
// @Failure      409 {object} dtos.apiError "The greeting already exists"
// @Failure      400 {object} dtos.apiError "The JSON payload was malformed"
// @Router       /api/v1/sample/greetings [post]
func (t SampleController) AddGreeting(ctx echo.Context, newGreeting NewGreetingRequest) error {
	requestCtx := request.ExtractContext(ctx)
	addErr := database.WithTransaction(requestCtx, func(dbCtx context.Context) error {
		return t.sampleLogic.AddGreeting(dbCtx, newGreeting.Greeting, t.greetingReader, t.greetingWriter)
	})

	// Decide how to respond based on business logic error
	if errors.Is(addErr, sample.ErrGreetingAlreadyExists) {
		logger.Log.Warn("Incoming request already existed.", zap.Error(addErr), zap.String("newGreeting", newGreeting.Greeting))
		apiErr := response.Conflict(addErr)
		apiErr.Description = "The provided greeting already exists in the system."
		return apiErr.Respond(ctx)
	} else if addErr != nil {
		logger.Log.Error("Something went wrong when adding greeting.", zap.Error(addErr), zap.String("newGreeting", newGreeting.Greeting))
		return response.InternalServerError(addErr).Respond(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}
