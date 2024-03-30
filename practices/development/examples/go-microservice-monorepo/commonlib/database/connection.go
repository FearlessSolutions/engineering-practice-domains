package database

import (
	"fmt"
	"strconv"
	"time"

	"example.com/sample/commonlib/config"
	"example.com/sample/commonlib/config/sharedoptions"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// Config contains configuration options for the database.
type Config struct {
	// Username is the username used to connect to the database.
	Username string

	// Password is the password used when connecting to the database.
	Password string

	// Host is the hostname or IP address of the MariaDB server to connect to.
	Host string

	// Schema is the default schema that should be used for SQL queries.
	Schema string

	// OptionalSettings contains settings that aren't required. You can safely leave this as the zero value if you want.
	OptionalSettings OptionalSettings
}

// OptionalSettings are performance-based connection settings that can be tweaked if desired.
type OptionalSettings struct {
	// Port is the database port to connect to, if something other than 3306.
	Port *int
	// MaxOpenConnections sets a limit on the number of total open MySQL connections. It is set to 20 by default, and
	// should be larger than MaxIdleConnections.
	MaxOpenConnections *int
	// MaxIdleConnections sets a limit on the limit of idle SQL connections in the database pool. It is set to 5 by default,
	// and should be smaller than MaxOpenConnections.
	MaxIdleConnections *int
}

// Connect connects to the database using the provided configuration.
func Connect(config Config) (*sqlx.DB, error) {
	dbHost := config.Host
	if config.OptionalSettings.Port != nil {
		dbHost += fmt.Sprintf(":%v", *config.OptionalSettings.Port)
	}

	db, connectErr := sqlx.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v", config.Username, config.Password, dbHost, config.Schema))
	if connectErr != nil {
		return nil, fmt.Errorf("failed to connect to database with given credentials (user %v, host %v, schema %v): %w",
			config.Username, dbHost, config.Schema, connectErr)
	}

	// Setting max connection lifetime to 3 minutes, as less than 5 minutes is recommended by the driver: https://github.com/go-sql-driver/mysql#important-settings
	db.SetConnMaxLifetime(3 * time.Minute)

	if config.OptionalSettings.MaxIdleConnections != nil {
		db.SetMaxIdleConns(*config.OptionalSettings.MaxIdleConnections)
	} else {
		db.SetMaxIdleConns(5)
	}
	if config.OptionalSettings.MaxOpenConnections != nil {
		db.SetMaxOpenConns(*config.OptionalSettings.MaxIdleConnections)
	} else {
		db.SetMaxOpenConns(20)
	}

	return db, nil
}

// ConnectFromConfig reads database configuration options from the environment, namely those listed in sharedoptions.DBOptions,
// and constructs the database.
func ConnectFromConfig(registry config.Registry) (*sqlx.DB, error) {
	dbConfig := Config{
		Username: registry.GetRequired(sharedoptions.DBUser),
		Password: registry.GetRequired(sharedoptions.DBPassword),
		Host:     registry.GetRequired(sharedoptions.DBHostname),
		Schema:   registry.GetRequired(sharedoptions.DBSchema),
	}
	if value, isPresent := registry.Get(sharedoptions.DBMaxConnections); isPresent {
		if intValue, parseErr := strconv.Atoi(value); parseErr == nil {
			dbConfig.OptionalSettings.MaxOpenConnections = &intValue
		} else {
			return nil, fmt.Errorf("a non-number slipped past validation on max database connections option with value \"%v\": %w", value, parseErr)
		}
	}
	if value, isPresent := registry.Get(sharedoptions.DBMaxIdleConnections); isPresent {
		if intValue, parseErr := strconv.Atoi(value); parseErr == nil {
			dbConfig.OptionalSettings.MaxIdleConnections = &intValue
		} else {
			return nil, fmt.Errorf("a non-number slipped past validation on max idle database connections option with value \"%v\": %w", value, parseErr)
		}
	}
	if value, isPresent := registry.Get(sharedoptions.DBPort); isPresent {
		if intValue, parseErr := strconv.Atoi(value); parseErr == nil {
			dbConfig.OptionalSettings.Port = &intValue
		} else {
			return nil, fmt.Errorf("a non-number slipped past validation on database port option with value \"%v\": %w", value, parseErr)
		}
	}

	return Connect(dbConfig)
}
