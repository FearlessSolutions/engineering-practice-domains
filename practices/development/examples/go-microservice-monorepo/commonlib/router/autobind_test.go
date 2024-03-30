package router

import (
	"net/http"
	"testing"

	"example.com/sample/commonlib/logger"
	reqhelper "example.com/sample/commonlib/request/testhelper"
	reshelper "example.com/sample/commonlib/response/testhelper"
	"github.com/jellydator/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
)

type AutoBindSuite struct {
	suite.Suite
}

func TestAutoBindSuite(t *testing.T) {
	suite.Run(t, new(AutoBindSuite))
}

func (suite *AutoBindSuite) SetupSuite() {
	initErr := logger.InitLogger(zapcore.InfoLevel, false)
	suite.Require().NoError(initErr)
}

type DummyRequestBody struct {
	SomeValue string
}

func (rb DummyRequestBody) Validate() error {
	return validation.ValidateStruct(&rb,
		validation.Field(&rb.SomeValue, validation.Required),
	)
}

type DummyResponseBody struct {
	PassedValue string `json:"passedValue"`
}

func DummyHTTPRoute(ctx echo.Context, requestBody DummyRequestBody) error {
	return ctx.JSON(http.StatusOK, DummyResponseBody{
		PassedValue: requestBody.SomeValue,
	})
}

func (suite *AutoBindSuite) TestAutoBindCallsNestedHandler() {
	var invoked bool
	body := DummyRequestBody{SomeValue: "hello"}

	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.GET, "/api/some/endpoint").
		WithBody(body).
		Build()
	suite.Require().NoError(buildErr)

	route := AutoBindAndValidate(func(ctx echo.Context, requestBody DummyRequestBody) error {
		invoked = true
		return DummyHTTPRoute(ctx, requestBody)
	})
	routeErr := route(request)
	suite.Require().NoError(routeErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusOK, response.StatusCode)
	suite.Require().True(invoked)

	var responseBody DummyResponseBody
	reshelper.UnmarshalBody(&suite.Suite, response.Body, &responseBody)
	suite.Require().Equal(body.SomeValue, responseBody.PassedValue)
}

func (suite *AutoBindSuite) TestAutoBindRespondsBadRequestOnBadJson() {
	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.GET, "/api/some/endpoint").
		WithRawBody([]byte("{bad json")).
		Build()
	suite.Require().NoError(buildErr)

	route := AutoBindAndValidate(DummyHTTPRoute)
	routeErr := route(request)
	suite.Require().NoError(routeErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusBadRequest, response.StatusCode)
}

func (suite *AutoBindSuite) TestAutoBindRespondsBadRequestOnInvalidBody() {
	bodyMissingRequiredFields := DummyRequestBody{}
	request, responseRecorder, buildErr := reqhelper.NewRequest(echo.GET, "/api/some/endpoint").
		WithBody(bodyMissingRequiredFields).
		Build()
	suite.Require().NoError(buildErr)

	route := AutoBindAndValidate(DummyHTTPRoute)
	routeErr := route(request)
	suite.Require().NoError(routeErr)

	response := responseRecorder.Result()
	suite.Require().Equal(http.StatusBadRequest, response.StatusCode)
}
