package sharedoptions

import (
	"example.com/sample/commonlib/config"
	"github.com/jellydator/validation"
	"github.com/jellydator/validation/is"
)

// IsInProduction is either "true" or "false" and determines whether or not the app is executing in production mode
var IsInProduction = config.NewValidatedOption("IS_PRODUCTION", true, func(value string) error {
	return validation.Validate(
		value,
		validation.In("true", "false").Error("value must be 'true' or 'false'"),
	)
})

// LogLevel determines the log level when starting the application. The value must be a valid Zap log level.
var LogLevel = config.NewValidatedOption("LOG_LEVEL", false, func(value string) error {
	return validation.Validate(
		value,
		validation.In("debug", "info", "warn", "error", "panic", "fatal").
			Error("value must be one of debug, info, warn, error, panic, or fatal"),
	)
})

// ListenPort determines the port the application listens on
var ListenPort = config.NewValidatedOption("LISTEN_PORT", false, func(value string) error {
	return validation.Validate(value, is.Port)
})

// AllowedOrigins contains a comma-separated list of allowed CORS origins
var AllowedOrigins = config.NewOption("ALLOWED_CORS_ORIGINS", false)

// DBUser is the username used to authenticate with the database
var DBUser = config.NewOption("DB_USER", true)

// DBPassword is the password used to authenticate with the database
var DBPassword = config.NewOption("DB_PASSWORD", true)

// DBHostname is the hostname of the database to connect to
var DBHostname = config.NewValidatedOption("DB_HOST", true, func(value string) error {
	return validation.Validate(value, is.Host)
})

// DBPort is the database port the application should connect to
var DBPort = config.NewValidatedOption("DB_PORT", false, func(value string) error {
	return validation.Validate(value, is.Int)
})

// DBSchema is the schema to use by default once connected to the database
var DBSchema = config.NewOption("DB_SCHEMA", true)

// DBMaxConnections is the number of total SQL connections the database pool cannot exceed
var DBMaxConnections = config.NewValidatedOption("DB_MAX_CONNECTIONS", false, func(value string) error {
	return validation.Validate(value, is.Int)
})

// DBMaxIdleConnections is the number of total idle SQL connections the database pool cannot exceed. This number should be less than DBMaxConnections.
var DBMaxIdleConnections = config.NewValidatedOption("DB_MAX_IDLE_CONNECTIONS", false, func(value string) error {
	return validation.Validate(value, is.Int.Error("should be a valid port number"))
})

// DBOptions is a bundle of all available database configuration options
var DBOptions = []config.Option{DBUser, DBPassword, DBHostname, DBPort, DBSchema, DBMaxConnections, DBMaxIdleConnections}
