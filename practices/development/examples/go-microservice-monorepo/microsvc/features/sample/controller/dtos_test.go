package controller

import (
	"github.com/jellydator/validation"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DtoSuite struct {
	suite.Suite
}

func TestDtoSuite(t *testing.T) {
	suite.Run(t, new(DtoSuite))
}

func (suite *DtoSuite) TestSampleGreetingValidationRequiredFields() {
	dtoWithNoRequiredFields := SampleGreetingRequest{}
	validationErr := dtoWithNoRequiredFields.Validate()
	suite.Require().Error(validationErr)

	var fieldErrors validation.Errors
	suite.Require().ErrorAs(validationErr, &fieldErrors)
	suite.Require().Contains(fieldErrors, "name")

	var nameError validation.ErrorObject
	suite.Require().ErrorAs(fieldErrors["name"], &nameError)
	suite.Require().Equal(validation.ErrRequired.Code(), nameError.Code())
}

func (suite *DtoSuite) TestNewGreetingValidationRequiredFields() {
	dtoWithNoRequiredFields := NewGreetingRequest{}
	validationErr := dtoWithNoRequiredFields.Validate()
	suite.Require().Error(validationErr)

	var fieldErrors validation.Errors
	suite.Require().ErrorAs(validationErr, &fieldErrors)
	suite.Require().Contains(fieldErrors, "greeting")

	var greetingError validation.ErrorObject
	suite.Require().ErrorAs(fieldErrors["greeting"], &greetingError)
	suite.Require().Equal(validation.ErrRequired.Code(), greetingError.Code())
}

func (suite *DtoSuite) TestNewGreetingValidationGreetingLength() {
	dtoWithTooLongGreting := NewGreetingRequest{
		Greeting: "This greeting is too long for the database to store oh noooooooo",
	}
	validationErr := dtoWithTooLongGreting.Validate()
	suite.Require().Error(validationErr)

	var fieldErrors validation.Errors
	suite.Require().ErrorAs(validationErr, &fieldErrors)
	suite.Require().Contains(fieldErrors, "greeting")

	var greetingError validation.ErrorObject
	suite.Require().ErrorAs(fieldErrors["greeting"], &greetingError)
	suite.Require().Equal(validation.ErrLengthTooLong.Code(), greetingError.Code())
}
