package config

import (
	"testing"

	"github.com/jellydator/validation"
	"github.com/stretchr/testify/suite"
)

type RegistrySuite struct {
	suite.Suite
}

func TestRegistrySuite(t *testing.T) {
	suite.Run(t, new(RegistrySuite))
}

func (suite *RegistrySuite) TestRejectsDuplicateOptions() {
	option := NewOption("MY_OPTION", true)
	builder := NewMockRegistryBuilder(map[string]string{
		option.VariableName(): "abcde",
	})
	builder.AddOption(option)
	suite.Require().Panics(func() {
		builder.AddOption(option)
	})
}

func (suite *RegistrySuite) TestPresenceVerification() {
	requiredOption := NewOption("MY_OPTION", true)
	notRequiredOption := NewOption("OTHER_OPTION", false)
	builder := NewMockRegistryBuilder(nil)
	builder.AddOptions([]Option{
		requiredOption,
		notRequiredOption,
	})

	_, buildErr := builder.VerifyAndBuild()

	suite.Require().Error(buildErr)

	var configurationError ErrIncorrectConfiguration
	suite.Require().ErrorAs(buildErr, &configurationError)
	suite.Require().Len(configurationError.MissingRequiredVariables, 1)
	suite.Require().Equal(requiredOption.VariableName(), configurationError.MissingRequiredVariables[0])
}

func (suite *RegistrySuite) TestValidationAppliesForPresentVariables() {
	validatedRequiredOption := NewValidatedOption("MUST_BE_CHOCOLATE", true, func(value string) error {
		return validation.Validate(value, validation.In("chocolate"))
	})
	validatedOptionalOption := NewValidatedOption("MUST_BE_SPICY", false, func(value string) error {
		return validation.Validate(value, validation.In("habanero", "jalapeno"))
	})
	validatedNotPresentRequiredOption := NewValidatedOption("NOT_HERE", true, func(value string) error {
		return validation.Validate(value, validation.Length(1, 10))
	})
	validatedNotPresentOptionalOption := NewValidatedOption("ALSO_NOT_HERE", false, func(value string) error {
		return validation.Validate(value, validation.Length(1, 10))
	})

	builder := NewMockRegistryBuilder(map[string]string{
		validatedRequiredOption.VariableName(): "vanilla",
		validatedOptionalOption.VariableName(): "mayonnaise",
	})
	builder.AddOptions([]Option{
		validatedRequiredOption,
		validatedOptionalOption,
		validatedNotPresentRequiredOption,
		validatedNotPresentOptionalOption,
	})
	_, buildErr := builder.VerifyAndBuild()

	suite.Require().Error(buildErr)

	var configErr ErrIncorrectConfiguration
	suite.Require().ErrorAs(buildErr, &configErr)

	// The missing required variable should register as missing
	suite.Require().Len(configErr.MissingRequiredVariables, 1)
	suite.Require().Equal(validatedNotPresentRequiredOption.VariableName(), configErr.MissingRequiredVariables[0])

	// Even though everything is validated, only the present variables should get validated regardless of whether they're required or not
	suite.Require().Len(configErr.InvalidVariables, 2)
	invalidVariableNames := []string{configErr.InvalidVariables[0].Name, configErr.InvalidVariables[1].Name}
	suite.Assert().Contains(invalidVariableNames, validatedRequiredOption.VariableName())
	suite.Assert().Contains(invalidVariableNames, validatedOptionalOption.VariableName())
}

func (suite *RegistrySuite) TestOptionPresence() {
	presentVariable := NewOption("I_AM_HERE", false)
	notPresentVariable := NewOption("I_AM_NOT_HERE", false)

	builder := NewMockRegistryBuilder(map[string]string{
		presentVariable.VariableName(): "hello",
	})
	builder.AddOptions([]Option{
		presentVariable,
		notPresentVariable,
	})
	registry, buildErr := builder.VerifyAndBuild()
	suite.Require().NoError(buildErr)

	presentValue, presentPresence := registry.Get(presentVariable)
	_, notPresentPresence := registry.Get(notPresentVariable)

	suite.Assert().True(presentPresence)
	suite.Assert().Equal("hello", presentValue)
	suite.Assert().False(notPresentPresence)
}

func (suite *RegistrySuite) TestUsingNonRegisteredOptionPanics() {
	option := NewOption("NOT_REGISTERED", false)
	builder := NewMockRegistryBuilder(nil)
	registry, buildErr := builder.VerifyAndBuild()

	suite.Require().NoError(buildErr)
	suite.Require().Panics(func() {
		_, _ = registry.Get(option)
	})
}

func (suite *RegistrySuite) TestRetrievingRequiredWithOptionalPanics() {
	optionalOption := NewOption("MY_OPTION", false)
	builder := NewMockRegistryBuilder(map[string]string{
		optionalOption.VariableName(): "abcde",
	})
	builder.AddOption(optionalOption)
	registry, buildErr := builder.VerifyAndBuild()

	suite.Require().NoError(buildErr)
	suite.Require().Panics(func() {
		_ = registry.GetRequired(optionalOption)
	})
}

func (suite *RegistrySuite) TestOptionsRetrieveSuccessfully() {
	optionalOption := NewOption("MY_OPTION", false)
	requiredOption := NewOption("OTHER_OPTION", true)
	builder := NewMockRegistryBuilder(map[string]string{
		optionalOption.VariableName(): "abcde",
		requiredOption.VariableName(): "wxyz",
	})
	builder.AddOptions([]Option{
		requiredOption,
		optionalOption,
	})
	registry, buildErr := builder.VerifyAndBuild()

	suite.Require().NoError(buildErr)

	optionalValue, optionalPresent := registry.Get(optionalOption)
	requiredValue := registry.GetRequired(requiredOption)

	suite.Assert().True(optionalPresent)
	suite.Assert().Equal("abcde", optionalValue)
	suite.Assert().Equal("wxyz", requiredValue)
}
