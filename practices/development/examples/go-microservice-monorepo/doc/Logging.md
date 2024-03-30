# Logging

Logging is achieved with a global logger, `logger.Log`. It supports logging at different levels, filtering out certain
log levels, and logging in several formats. The logger is from Uber's [Zap](https://github.com/uber-go/zap) logging library.

## The global logger instance

Any piece of code can log via the global logger instance. Additional information can be logged with typed `zap.Field`s 
which can be constructed via functions listed in [this section of the zap documentation](https://pkg.go.dev/go.uber.org/zap#Field).

Here's an example of logging with additional related data:

```go
// Log that you fell asleep counting sheep, along with the number of sheep counted and some of the sheep's names
sheepCounted := 218
sheepNames := []string{"Wooly", "Seamus", "Dave"}
// Additional data accepts key/value pairs
logger.Log.Info("Fell asleep by counting sheep!", zap.Int("sheepCounted", sheepCounted), zap.Strings("sheepNames", sheepNames))

// Log that an error occurred
someError := errors.New("trapped in a dream")
// Errors can be added with zap.Error, which implicitly has the key "error". zap.NamedError can log an error with a different key.
logger.Log.Error("Failed to wake up again! Oh no!", zap.Error(someError))
```

### Log level recommendations

Being able to filter log levels is only useful when you can filter out certain sets of data. Here's the sort of information
that is recommended at certain log levels:

* **debug** - Should be used for informational logging describing implementation details of various layers within the application. It's useful, but can be noisy
* **info** - Should be used to describe **what** operations are occurring inside the app instead of **why**. Typical use cases are logging that an endpoint was invoked and how the application responded.
* **warn** - Should be used to describe **expected errors** such as error cases explicitly defined by the business logic or validation failures.
* **error** - Should be used to describe **unexpected errors** such as bad JSON or database connectivity errors that are handled by catch-all error handlers.
* **panic** - Should not be used, especially as panics will just be caught by the Recovery middleware (see [Middleware.md](Middleware.md#recovery-middleware) for more information).
* **fatal** - Should be used for logging the reason the application is in an invalid state and cannot start up. Logging at this level calls `os.Exit()`, terminating the program.


## Production and dev mode

The global logger has two output formats depending on whether the application is running in production or not. In development, 
the logger uses a plaintext format easily understood by humans. In production, the logger logs in a JSON format that can
easily be aggregated and read by log aggregation services such as Grafana. If the logger is instantiated based on the application
configuration, the "production" state is determined by the `sharedoptions.IsInProduction` configuration option. Otherwise, it's a
simple constructor parameter. 

See ["Instantiating the logger"](#instantiating-the-logger) for more information on constructing the
logger, and [the repository layout docs](Navigation and Repository Layout.md#layout-of-the-common-library) 
for more information on shared configuration options.

## Adjusting the log level

Log level filters are controlled by both the shared option `sharedoptions.LogLevel` and the shared log level controller defined in the `loglevel` shared feature.
See [the repository layout docs](Navigation and Repository Layout.md) for where to find those packages in the repo.

The `loglevel` shared feature exposes HTTP endpoints for live-updating the application log level in production without an application restart.

If you need to write code to adjust the global log level, you can invoke the function `logger.AdjustLevel()` to change the log level,
passing a valid [zapcore.Level](https://pkg.go.dev/go.uber.org/zap/zapcore@v1.26.0#Level) to adjust the log level. Here's an example:

```go
// Set the global log level to "debug"
logger.AdjustLevel(zapcore.DebugLevel)
```

## Instantiating the logger

You may need to instantiate the logger in various situations, such as creating a new microservice or instantiating the logger
during tests. There are two ways to do so:

1. Creating the logger based on application configuration with a `config.Registry`, or
2. Creating the logger directly via its constructor

For more information on configuration registries, see [the configuration docs](Configuration.md#creating-a-configuration-registry-and-registering-options).

### Instantiating the logger directly

To initialize the global logger, you can call `logger.InitLogger()`, passing it its initial log level and a boolean
representing whether to output in production mode or not. Here's an example:

```go
// Initialize the logger at "warn" level in development mode:
initializationErr := logger.InitLogger(zapcore.LevelWarn, false)
if initializationErr != nil {
    fmt.Println("Initializing the logger failed!")
}

// Initialize the logger at "info" level in production mode:
initializationErr := logger.InitLogger(zapcore.LevelInfo, true)
if initializationErr != nil {
    fmt.Println("Initializing the logger failed!")
}
```

### Instantiating the logger from a configuration registry

To initialize the global logger from a configuration registry, consult the `logger.InitLoggerFromConfig` inline docs for what
shared options you'll need registered with the registry. You can find out how to construct a registry in [the configuration docs](Configuration.md#creating-a-configuration-registry-and-registering-options).
With your constructed registry, you can call `logger.InitLoggerFromConfig()`, passing it the registry. Here's an example,
using a microservice's global registry in the `options` package:

```go
// Construct the logger based on config options
initializationErr := logger.InitLoggerFromConfig(*options.Registry)
if initializationErr != nil {
	fmt.Println("Initializing the logger failed!")
}
```