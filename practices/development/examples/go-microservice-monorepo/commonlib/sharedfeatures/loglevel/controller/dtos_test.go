package controller

import (
	"testing"

	"github.com/jellydator/validation"
	"github.com/stretchr/testify/suite"
)

type DtosSuite struct {
	suite.Suite
}

func TestDtosSuite(t *testing.T) {
	suite.Run(t, new(DtosSuite))
}

func (suite *DtosSuite) TestLogLevelRequestRequiredFields() {
	request := ChangeLogLevelRequest{}
	validationErr := request.Validate()
	suite.Require().Error(validationErr)

	var fieldErrors validation.Errors
	suite.Require().ErrorAs(validationErr, &fieldErrors)
	suite.Require().Contains(fieldErrors, "NewLevel")

	var errorInfo validation.ErrorObject
	suite.Require().ErrorAs(fieldErrors["NewLevel"], &errorInfo)
	suite.Require().Equal(validation.ErrRequired.Code(), errorInfo.Code())
}

func (suite *DtosSuite) TestLogLevelRequestOnlyAcceptsLogLevels() {
	type levelTest struct {
		testName             string
		levelInput           string
		shouldPassValidation bool
	}

	testCases := []levelTest{
		{
			testName:             "Should accept debug level",
			levelInput:           "debug",
			shouldPassValidation: true,
		},
		{
			testName:             "Should accept info level",
			levelInput:           "info",
			shouldPassValidation: true,
		},
		{
			testName:             "Should accept warn level",
			levelInput:           "warn",
			shouldPassValidation: true,
		},
		{
			testName:             "Should accept error level",
			levelInput:           "error",
			shouldPassValidation: true,
		},
		{
			testName:             "Should accept panic level",
			levelInput:           "panic",
			shouldPassValidation: true,
		},
		{
			testName:             "Should accept fatal level",
			levelInput:           "fatal",
			shouldPassValidation: true,
		},
		{
			testName:             "Should reject non-levels",
			levelInput:           "not a log level",
			shouldPassValidation: false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.testName, func() {
			request := ChangeLogLevelRequest{NewLevel: testCase.levelInput}
			validationErr := request.Validate()
			if testCase.shouldPassValidation {
				suite.Require().NoError(validationErr)
			} else {
				suite.Require().Error(validationErr)
			}
		})
	}
}
