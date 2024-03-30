package config

import (
	"fmt"
	"strings"
)

// RegistryBuilder is a builder for a Registry. It accepts the full list of configuration options available to
// the app and verifies they look OK during construction of the registry.
type RegistryBuilder struct {
	registeredOptions map[string]struct{}
	allOptions        []Option
	dataSource        registrySource
}

// NewRegistryBuilder constructs a RegistryBuilder
func NewRegistryBuilder() RegistryBuilder {
	return RegistryBuilder{
		registeredOptions: make(map[string]struct{}),
		allOptions:        nil,
		dataSource:        environmentRegistrySource{},
	}
}

// NewMockRegistryBuilder constructs a RegistryBuilder that uses the passed map of environment variables in place of
// the real environment. Nil can be passed for an empty environment
func NewMockRegistryBuilder(environmentVariables map[string]string) RegistryBuilder {
	if environmentVariables == nil {
		environmentVariables = make(map[string]string)
	}
	return RegistryBuilder{
		registeredOptions: make(map[string]struct{}),
		allOptions:        nil,
		dataSource:        mockRegistrySource(environmentVariables),
	}
}

// AddOption registers a new Option with the builder
func (rb *RegistryBuilder) AddOption(option Option) {
	if _, optionExists := rb.registeredOptions[option.envName]; optionExists {
		panic(fmt.Sprintf("Registering option failed, nother option with the same name exists: %v", option.envName))
	}

	rb.registeredOptions[option.envName] = struct{}{}
	rb.allOptions = append(rb.allOptions, option)
}

// AddOptions registers a list of Options with the builder
func (rb *RegistryBuilder) AddOptions(options []Option) {
	for _, option := range options {
		rb.AddOption(option)
	}
}

// InvalidVariable describes a validation error with an environment variable
type InvalidVariable struct {
	Name            string
	ValidationError error
}

// ErrIncorrectConfiguration is an error which describes issues with environment configuration, either for
// required variables or variables that failed validation
type ErrIncorrectConfiguration struct {
	MissingRequiredVariables []string
	InvalidVariables         []InvalidVariable
}

// errorsPresent returns true if any required variables are missing or variables failed validation
func (err ErrIncorrectConfiguration) errorsPresent() bool {
	return len(err.MissingRequiredVariables) > 0 || len(err.InvalidVariables) > 0
}

// Error implements the error interface for ErrIncorrectConfiguration
func (err ErrIncorrectConfiguration) Error() string {
	var errorText string
	if len(err.MissingRequiredVariables) > 0 {
		variablesMissingList := strings.Join(err.MissingRequiredVariables, ", ")
		errorText += fmt.Sprintf("the following required environment variables were missing: %v", variablesMissingList)
	}

	if len(err.InvalidVariables) > 0 {
		var validationErrorsList []string
		for _, validationError := range err.InvalidVariables {
			validationErrorsList = append(validationErrorsList, fmt.Sprintf("%v was invalid for this reason: %v", validationError.Name, validationError.ValidationError.Error()))
		}
		validationErrors := strings.Join(validationErrorsList, ", ")

		if len(errorText) == 0 {
			errorText = validationErrors
		} else {
			errorText += ", " + validationErrors
		}
	}

	return errorText
}

// Is allows someone to verify an error has a nested instance of ErrIncorrectConfiguration using errors.Is
//
//goland:noinspection GoTypeAssertionOnErrors
func (err ErrIncorrectConfiguration) Is(testErr error) bool {
	if _, isErrType := testErr.(ErrIncorrectConfiguration); isErrType {
		return true
	} else if _, isErrType = testErr.(*ErrIncorrectConfiguration); isErrType {
		return true
	}

	return false
}

// VerifyAndBuild checks the registered Options against the environment, producing a constructed registry if
// all required environment variables are present and all environment variables with validation pass validation.
// Produces an ErrIncorrectConfiguration describing problems with the current configuration if those checks fail.
func (rb *RegistryBuilder) VerifyAndBuild() (Registry, error) {
	var variableIssues ErrIncorrectConfiguration
	for _, option := range rb.allOptions {
		optionValue, optionPresent := rb.dataSource.getValue(option.envName)
		// Verify required option is present
		if option.required {
			if !optionPresent {
				variableIssues.MissingRequiredVariables = append(variableIssues.MissingRequiredVariables, option.envName)
			}
		}

		// If the option is present at all, and it has a validation function, validate the option
		if optionPresent && option.validationFn != nil {
			validationErr := option.validationFn(optionValue)
			if validationErr != nil {
				variableIssues.InvalidVariables = append(variableIssues.InvalidVariables, InvalidVariable{
					Name:            option.envName,
					ValidationError: validationErr,
				})
			}
		}
	}

	if variableIssues.errorsPresent() {
		return Registry{}, variableIssues
	}

	return Registry{
		registeredOptions: rb.registeredOptions,
		dataSource:        rb.dataSource,
	}, nil
}

// Registry is a validated registry of Option values. This type should be constructed via a RegistryBuilder which will
// verify the integrity of environment variables during the construction of this type
type Registry struct {
	registeredOptions map[string]struct{}
	dataSource        registrySource
}

// Get retrieves the specified Option from the environment. The first return value contains the option's value if it's
// present, and the second value is true if the value was actually present or false otherwise. This function will panic
// if an unregistered option is passed.
func (reg Registry) Get(option Option) (string, bool) {
	if _, registeredOption := reg.registeredOptions[option.envName]; !registeredOption {
		panic(fmt.Sprintf("Option %v is not registered! Make sure to register it when building the option registry.", option.envName))
	}

	return reg.dataSource.getValue(option.envName)
}

// GetRequired retrieves the specified required Option from the environment. Since the passed option is required,
// the "presence" boolean is not returned like it is in Get. This function will panic if it is passed an optional
// option or an unregistered option.
func (reg Registry) GetRequired(option Option) string {
	if !option.required {
		panic(fmt.Sprintf("Non-required option %v was passed to Registry.GetRequired", option.envName))
	}
	// We don't need to check presence since the construction of the registry verified the option is present
	value, _ := reg.Get(option)
	return value
}
