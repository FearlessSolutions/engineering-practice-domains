package options

import (
	"example.com/sample/commonlib/config"
	"example.com/sample/commonlib/config/sharedoptions"
)

// Registry is the global configuration registry of verified environment variables
var Registry *config.Registry

// buildRegistry attaches configuration options to the passed config.RegistryBuilder and attempts the verification
// and build process to verify environment variables are what we expect. It returns an error if required variables
// aren't present or some environment variables didn't pass validation
func buildRegistry(regBuilder config.RegistryBuilder) error {
	regBuilder.AddOptions([]config.Option{
		sharedoptions.IsInProduction,
		sharedoptions.LogLevel,
		sharedoptions.AllowedOrigins,
		sharedoptions.ListenPort,
	})
	regBuilder.AddOptions(sharedoptions.DBOptions)

	registry, buildErr := regBuilder.VerifyAndBuild()
	if buildErr != nil {
		return buildErr
	}

	Registry = &registry
	return nil
}

// InitRegistry initializes the global configuration registry, Registry
func InitRegistry() error {
	regBuilder := config.NewRegistryBuilder()
	return buildRegistry(regBuilder)
}
