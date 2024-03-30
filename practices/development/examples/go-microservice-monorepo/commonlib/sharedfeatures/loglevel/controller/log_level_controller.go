package controller

import (
	"net/http"

	"example.com/sample/commonlib/logger"
	"example.com/sample/commonlib/response"
	"example.com/sample/commonlib/router"
	"example.com/sample/commonlib/sharedfeatures/loglevel"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LogLevelController is a REST controller that allows live updates of the application log level.
type LogLevelController struct {
	logicCore loglevel.Core
}

// New constructs a new LogLevelController
func New() LogLevelController {
	return LogLevelController{
		logicCore: loglevel.CoreLogic{},
	}
}

// newWithCore constructs a LogLevelController with a mocked core implementation
func newWithCore(core loglevel.Core) LogLevelController {
	return LogLevelController{
		logicCore: core,
	}
}

// AttachRoutes implements router.Controller for LogLevelController. It defines this controller's routes
func (ctrl LogLevelController) AttachRoutes(rtr *echo.Echo) {
	rtr.POST("/api/v1/config/log-level", router.AutoBindAndValidate(ctrl.AdjustLogLevel))
}

// AdjustLogLevel is a route that triggers business logic to adjust the app's log level
func (ctrl LogLevelController) AdjustLogLevel(ctx echo.Context, changeReq ChangeLogLevelRequest) error {
	level, parseErr := zap.ParseAtomicLevel(changeReq.NewLevel)
	if parseErr != nil {
		logger.Log.Error("Should have been able to parse log level based on validation, but it failed.", zap.Error(parseErr))
		return response.InternalServerError(parseErr).Respond(ctx)
	}

	ctrl.logicCore.SetLogLevel(level.Level())

	return ctx.NoContent(http.StatusOK)
}
