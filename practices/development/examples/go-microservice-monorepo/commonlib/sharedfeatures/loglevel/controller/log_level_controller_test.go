package controller

import (
	"net/http"
	"testing"

	"example.com/sample/commonlib/logger"
	reqhelper "example.com/sample/commonlib/request/testhelper"
	"example.com/sample/commonlib/router"
	"example.com/sample/commonlib/sharedfeatures/loglevel"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zapcore"
)

type LogLevelControllerSuite struct {
	suite.Suite
	mockController *gomock.Controller
	coreMock       *loglevel.MockCore
}

func TestLogLevelControllerSuite(t *testing.T) {
	suite.Run(t, new(LogLevelControllerSuite))
}

func (suite *LogLevelControllerSuite) SetupSuite() {
	setupErr := logger.InitLogger(zapcore.DebugLevel, false)
	suite.Require().NoError(setupErr)
}

func (suite *LogLevelControllerSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.coreMock = loglevel.NewMockCore(suite.mockController)
}

func (suite *LogLevelControllerSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *LogLevelControllerSuite) TestLogLevelCanBeAdjusted() {
	suite.coreMock.EXPECT().SetLogLevel(zapcore.DebugLevel)

	requestBody := ChangeLogLevelRequest{
		NewLevel: "debug",
	}
	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.POST, "/api/v1/config/log-level").
		WithBody(requestBody).
		Build()
	suite.Require().NoError(buildErr)

	ctrl := newWithCore(suite.coreMock)
	route := router.AutoBindAndValidate(ctrl.AdjustLogLevel)
	responseErr := route(request)
	suite.Require().NoError(responseErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusOK, response.StatusCode)
}
