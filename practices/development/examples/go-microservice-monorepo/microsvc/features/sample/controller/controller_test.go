package controller

import (
	"example.com/sample/commonlib/logger"
	reqhelper "example.com/sample/commonlib/request/testhelper"
	reshelper "example.com/sample/commonlib/response/testhelper"
	"example.com/sample/commonlib/router"
	"example.com/sample/microsvc/features/sample"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zapcore"
	"net/http"
	"testing"
)

type SampleControllerSuite struct {
	suite.Suite
	mockController *gomock.Controller
	mockCore       *sample.MockCore
}

func TestSampleControllerSuite(t *testing.T) {
	suite.Run(t, new(SampleControllerSuite))
}

func (suite *SampleControllerSuite) SetupSuite() {
	logInitErr := logger.InitLogger(zapcore.InfoLevel, false)
	suite.Require().NoError(logInitErr)
}

func (suite *SampleControllerSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockCore = sample.NewMockCore(suite.mockController)
}

func (suite *SampleControllerSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *SampleControllerSuite) TestProduceGreeting() {
	suite.mockCore.EXPECT().GiveGreeting(gomock.Any(), "Evan", gomock.Any()).Return("Hello, Evan!", nil)
	reqBody := SampleGreetingRequest{Name: "Evan"}
	ctx, responseRecorder, buildErr := reqhelper.NewRequest(echo.POST, "/api/v1/sample/greeting").
		WithBody(reqBody).
		Build()
	suite.Require().NoError(buildErr)

	ctrl := newWithCore(suite.mockCore)
	route := router.AutoBindAndValidate(ctrl.ProduceGreeting)
	responseErr := route(ctx)
	suite.Require().NoError(responseErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusOK, response.StatusCode)

	var responseBody SampleGreetingResponse
	reshelper.UnmarshalBody(&suite.Suite, response.Body, &responseBody)

	suite.Assert().Equal(
		SampleGreetingResponse{
			Greeting: "Hello, Evan!",
		},
		responseBody,
	)
}

func (suite *SampleControllerSuite) TestAddGreeting() {
	suite.mockCore.EXPECT().AddGreeting(gomock.Any(), "G'day", gomock.Any(), gomock.Any()).Return(nil)
	requestBody := NewGreetingRequest{
		Greeting: "G'day",
	}
	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.POST, "/api/v1/sample/greeting/add").
		WithBody(requestBody).
		Build()
	suite.Require().NoError(buildErr)

	ctrl := newWithCore(suite.mockCore)
	route := router.AutoBindAndValidate(ctrl.AddGreeting)
	responseErr := route(request)
	suite.Require().NoError(responseErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusCreated, response.StatusCode)
}

func (suite *SampleControllerSuite) TestAddGreetingRespondsConflictOnExistingGreeting() {
	suite.mockCore.EXPECT().
		AddGreeting(gomock.Any(), "G'day", gomock.Any(), gomock.Any()).
		Return(sample.ErrGreetingAlreadyExists)
	requestBody := NewGreetingRequest{
		Greeting: "G'day",
	}
	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.POST, "/api/v1/sample/greeting/add").
		Build()
	suite.Require().NoError(buildErr)

	ctrl := newWithCore(suite.mockCore)
	responseErr := ctrl.AddGreeting(request, requestBody)
	suite.Require().NoError(responseErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusConflict, response.StatusCode)
}
