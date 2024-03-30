package sharedoptions

import (
	"testing"

	"example.com/sample/commonlib/config"
	"github.com/stretchr/testify/suite"
)

type CommonOptionsSuite struct {
	suite.Suite
}

type validationSubtestParams struct {
	testName             string
	registryValue        string
	shouldPassValidation bool
}

func TestCommonOptionsSuite(t *testing.T) {
	suite.Run(t, new(CommonOptionsSuite))
}

func (suite *CommonOptionsSuite) TestIsInProductionIsRequired() {
	builder := config.NewMockRegistryBuilder(nil)
	builder.AddOption(IsInProduction)
	_, rawBuildErr := builder.VerifyAndBuild()
	suite.Require().Error(rawBuildErr)

	var buildError config.ErrIncorrectConfiguration
	suite.Require().ErrorAs(rawBuildErr, &buildError)
	suite.Require().Contains(buildError.MissingRequiredVariables, IsInProduction.VariableName())
}

func (suite *CommonOptionsSuite) TestIsInProductionValidation() {

	subtests := []validationSubtestParams{
		{
			testName:             "Fails on not true or false",
			registryValue:        "not true or false",
			shouldPassValidation: false,
		},
		{
			testName:             "Succeeds on true",
			registryValue:        "true",
			shouldPassValidation: true,
		},
		{
			testName:             "Succeeds on false",
			registryValue:        "false",
			shouldPassValidation: true,
		},
	}

	for _, subtest := range subtests {
		suite.Run(subtest.testName, func() {
			builder := config.NewMockRegistryBuilder(map[string]string{
				IsInProduction.VariableName(): subtest.registryValue,
			})
			builder.AddOption(IsInProduction)
			_, rawBuildErr := builder.VerifyAndBuild()

			if subtest.shouldPassValidation {
				suite.Require().NoError(rawBuildErr)
			} else {
				suite.Require().Error(rawBuildErr)

				var buildError config.ErrIncorrectConfiguration
				suite.Require().ErrorAs(rawBuildErr, &buildError)
				suite.Require().Len(buildError.InvalidVariables, 1)
				suite.Require().Equal(IsInProduction.VariableName(), buildError.InvalidVariables[0].Name)
			}
		})
	}
}

func (suite *CommonOptionsSuite) TestLogLevelIsNotRequired() {
	builder := config.NewMockRegistryBuilder(nil)
	builder.AddOption(LogLevel)
	_, buildErr := builder.VerifyAndBuild()

	suite.Require().NoError(buildErr)
}

func (suite *CommonOptionsSuite) TestLogLevelValidation() {
	subtests := []validationSubtestParams{
		{
			testName:             "Rejects bad log level",
			registryValue:        "achoo",
			shouldPassValidation: false,
		},
		{
			testName:             "Accepts debug level",
			registryValue:        "debug",
			shouldPassValidation: true,
		},
		{
			testName:             "Accepts info level",
			registryValue:        "info",
			shouldPassValidation: true,
		},
		{
			testName:             "Accepts warn level",
			registryValue:        "warn",
			shouldPassValidation: true,
		},
		{
			testName:             "Accepts error level",
			registryValue:        "error",
			shouldPassValidation: true,
		},
		{
			testName:             "Accepts panic level",
			registryValue:        "panic",
			shouldPassValidation: true,
		},
		{
			testName:             "Accepts fatal level",
			registryValue:        "fatal",
			shouldPassValidation: true,
		},
	}

	for _, subtest := range subtests {
		suite.Run(subtest.testName, func() {
			registryBuilder := config.NewMockRegistryBuilder(map[string]string{
				LogLevel.VariableName(): subtest.registryValue,
			})
			registryBuilder.AddOption(LogLevel)
			_, buildErr := registryBuilder.VerifyAndBuild()

			if subtest.shouldPassValidation {
				suite.Require().NoError(buildErr)
			} else {
				suite.Require().Error(buildErr)

				var validationError config.ErrIncorrectConfiguration
				suite.Require().ErrorAs(buildErr, &validationError)
				suite.Require().Len(validationError.InvalidVariables, 1)
				suite.Require().Equal(LogLevel.VariableName(), validationError.InvalidVariables[0].Name)
			}
		})
	}
}
