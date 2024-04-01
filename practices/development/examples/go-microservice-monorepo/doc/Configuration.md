# Configuration

The Go microservice template employs a sophisticated configuration registry system which can catch misconfiguration during
application startup and provide detailed diagnostics on why the application is misconfigured.

## Retrieving Configuration Options

Configuration options are retrieved from a global **configuration registry** constructed during startup. There are
two ways to retrieve configuration options for required or optional configuration options: `Registry.Get()` and `Registry.GetRequired()`.

### Registry.Get()

`Registry.Get()` retrieves either an optional or required configuration option from the environment. Here's an example of its
usage within a microservice:

```go
// The "options" package contains our global registry within a microservice
configurationValue, valuePresent := options.Registry.Get(options.ImportantConfigValue)
// Registry.Get follows the "comma ok" idiom described in Effective Go: https://go.dev/doc/effective_go#maps
if valuePresent {
	fmt.Printf("The configuration option is set! Here is its value: %v", configurationValue)
} else {
	fmt.Printf("The configuration option isn't present.")
}
```

Note that the configuration option must be registered with the configuration registry. See more information in [the section on options](#creating-and-registering-configuration-options).

### Registry.GetRequired()

`Registry.GetRequired()` retrieves a required configuration option from the environment. Note that the application will crash and report
an optional configuration option if it is passed to this function (this will help you catch these issues during development). It's nearly
identical to `Registry.Get()` except it doesn't return a boolean for the presence of the variable, and it more accurately expresses the intent
of extracting a required value.

Here's an example of its usage:

```go
// The "options" package contains our global registry with a microservice
configurationValue := options.Registry.GetRequired(options.SomeRequiredOption)
fmt.Printf("I'm guaranteed to have this value: %v", configurationValue)
```

## Creating and registering configuration options

The configuration registry exists because it cannot be constructed unless all registered options have been deemed valid, otherwise a 
diagnostic is reported. During construction of the registry, it checks the following:

1. All registered required options are present in the environment
2. All registered options that are present and have a registered validation function pass validation

If those checks do not pass, construction of the registry fails with an error describing the configuration problems.

### Creating a configuration option

Configuration options are created with the `config.NewOption()` function. It takes two parameters - the name of the configuration option
in the environment and whether the option should be required or not. These options are defined in the **options** package if
they're specific to a certain microservice or **sharedoptions** if the option is shared across multiple microservices. See [the repository layout docs](Navigation%20and%20Repository%20Layout.md)
for information on where those are. These options can also be documented, so one can see what a configuration option is for when
retrieved from a registry.

Here's a code example:

```go
// This is a required configuration option, defined outside any functions so that it can be used anywhere
var TheSecretIngredient = config.NewOption("SECRET_INGREDIENT", true)

// This is an optional configuration option, denoted by the second parameter
var TheAdditionalTopping = config.NewOption("ADDITIONAL_TOPPING", false)
```

### Creating a validated configuration option

Configuration options can also use `jellydator/validation` validators to validate the content of configuration options. See
the [validation documentation](Microservice%20Architecture.md#dtos-validation-and-responding) for more information on validation.
These validations are run during the construction of the configuration registry to verify everything is in the correct format.
Validations only run if the option is present, in the case of optional validation options.
To create a validated configuration option, you can either construct it via the alternate constructor `config.NewValidatedOption()`
or by calling `Option.SetValidation()` on an existing option. `config.NewValidatedOption()` is the same as `config.NewOption()` but it
takes a validation function in the 3rd parameter.

Here's a code example for a configuration option that has a maximum length:

```go
// ShortConfigOption is an optional configuration setting that can't be longer than 32 characters if present
var ShortConfigOption = config.NewValidatedOption("SHORT_CONFIG_OPTION", false, func(value string) error {
	validation.Validate(value, validation.Length(0, 32))
})
```

### Creating a configuration registry and registering options

Configuration options are registered with a configuration registry via a `config.RegistryBuilder` which should be defined in
a microservice's **options** package. See the [repository layout docs](Navigation%20and%20Repository%20Layout.md) for where that is.
Once all the options are added, the `RegistryBuilder.VerifyAndBuild()` function is invoked which verifies the configuration in the
environment and constructs the registry.

It should be noted that once the registry is initially created, **it cannot register more configuration options**. Options are added
one at a time with `Registry.AddOption()` or via a slice of them with `Registry.AddOptions()`.

Here is an example of creating a couple options, registering them with the registry, and constructing the registry:

```go
// Define config options
var RequiredOption = config.NewOption("I_AM_REQUIRED", true)
var ValidatedOption = config.NewValidatedOption("I_AM_A_NUMBER", false, func(value string) error {
	validation.Validate(value, is.Int)
})

func ConstructRegistry() Registry {
    // Get a registry builder
    registryBuilder := config.NewRegistryBuilder()
    // Register configuration options
    registryBulider.AddOptions([]config.Option{
        RequiredOption,
        ValidatedOption,
    })
	
    // Build the registry
    constructedRegistry, buildErr := registryBuilder.VerifyAndBuild()
    if buildErr != nil {
        errorMessage := fmt.Sprintf("Could not build the registry: %v", buildErr)
        // In this case, we'll just panic if the configuration is invalid to make this simpler
        panic(errorMessage)
    }
	
    return constructedRegistry
}
```
