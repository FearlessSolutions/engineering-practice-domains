package config

// Option represents an environment variable which can be verified during construction of a configuration Registry
// for its presence, or optionally validated by a provided function.
type Option struct {
	envName      string
	required     bool
	validationFn func(string) error
}

// NewOption constructs a new Option, accepting the name of the represented environment variable and whether
// it's required.
func NewOption(envName string, required bool) Option {
	return Option{
		envName:  envName,
		required: required,
	}
}

// NewValidatedOption constructs a new Option. This constructor works just like NewOption, but it accepts a validation
// function to verify the format of the represented environment variable.
func NewValidatedOption(envName string, required bool, validationFn func(value string) error) Option {
	return Option{
		envName:      envName,
		required:     required,
		validationFn: validationFn,
	}
}

// SetValidation adds a validation function to an Option.
func (opt *Option) SetValidation(validationFn func(value string) error) {
	opt.validationFn = validationFn
}

// VariableName returns the name of the environment variable this option represents
func (opt *Option) VariableName() string {
	return opt.envName
}
