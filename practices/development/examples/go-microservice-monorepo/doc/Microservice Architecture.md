# Architecture and writing different app layers

This microservice template is built using Hexagonal Architecture principles for maximum testability. [This article on medium](https://medium.com/ssense-tech/hexagonal-architecture-there-are-always-two-sides-to-every-story-bc0780ed7d9c)
provides a good overview of those principles - driving and driven ports, the business logic
core, and adapters. **It is highly recommended that you understand the architecture principles before continuing.**
This documentation heavily references different components of the architecture, so it will be good to understand its terms.

## Writing business logic

Business logic represents the core of the application. All business logic should be able to safely communicate with
other business logic. Code that triggers the business logic from outside the application, or code that accesses
external systems should be constrained to the adapters that plug into the ports the business logic defines.

A piece of business logic may define many ports - that is it can implement multiple interfaces used by external-facing
REST controllers and depend on multiple interfaces for retrieving data. There's not necessarily a one-to-one correlation
between each layer - it's all about separation of concerns.

Here's what a piece of business logic might look like, exposing itself on one driving port and depending on 2 driven ports:

<details>
<summary>Business logic code example</summary>

```go
package colorpalette

import (
	"context"
	"image/color"
)

// ColorReader is a driven port that allows the business logic to access a set of colors.
type ColorReader interface {
	// List lists the set of colors available
	List(ctx context.Context) ([]color.RGBA, error)
}

// ColorWriter is a driven port that allows the business logic to mutate the set of colors.
type ColorWriter interface {
	// Add adds a new color to the set of colors available
	Add(ctx context.Context, color color.RGBA)
}

// Core is the driving port for the business logic
type Core interface {
	// MixColors selects numToMix colors from the colors available and prints their mixed value
	MixColors(ctx context.Context, numToMix int, reader ColorReader) (color.RGBA, error)
}

// CoreLogic represents the business logic implementation. It implements the Core interface.
type CoreLogic struct{}

// MixColors implements Core for CoreLogic
func (CoreLogic) MixColors(ctx context.Context, numToMix int, reader ColorReader) (color.RGBA, error) {
	// ...Business logic here...
}
```
</details>


In this example, `CoreLogic` implements the driving port, `Core`. Note that the `MixColors` function accepts the driven
ports at the time the function is invoked - this is intended to allow the business logic to request specifically the
driven ports that it needs for that particular operation, and makes calling business logic across multiple features
easier, as the individual pieces of business logic don't need to be instantiated separately or nested inside one another.

This also helps expose the kind of capabilities a piece of business logic requires - if it accepted something like a
`PermissionsReader` then you can tell just by the call signature that the operation might be role-restricted. Details
on how the data is accessed should be confined to the driven port interface and should not be exposed to the business logic.

You may also notice that a `context.Context` is passed around - Go Contexts are used to cancel operations that are in progress
and can be used to propagate information through layers without exposing implementation detail. This technique should be
used only when incredibly necessary, but this is how database connections are transferred to database adapters and transactions
are propagated across business logic - see [the middleware docs](Middleware.md#database-connection-middleware) and
[the database adapter section](#acquiring-a-database-connection) for more information.

Note also that the business logic implementation just references the interface in its doc comment - this is intentional. Since driving adapters
communicate with the business logic over an interface, it makes sense to maintain the documentation for how the business logic
works on that interface rather than trying to maintain it in two places. Therefore, if calling into other business logic,
look for the referenced interface from the business logic's doc comment.

### Domain objects

Because we're using hexagonal architecture, the business logic should define data structures that are used to communicate with
it - adapters communicating through driving ports should convert their data types to domain objects before communicating with
the business logic, and driven ports should accept domain objects and perform any necessary conversions before sending
data to external systems. This allows the business logic to ignore implementation details and allows us to only expose fields
we want to expose to the outside.

### Errors in business logic

Since the business logic defines the contract for the ports it asks adapters to use or implement, the business logic
should also define the set of errors it expects to encounter. It should document those errors in documentation comments on the
port interface. Any other error returned by the adapter should be treated as an **adapter error**
and be immediately passed up to the caller. Here's an example of how errors might be exposed on various ports:

<details>
<summary>Domain error example</summary>

```go
// Package grouper provides calculations for managing groups of things
package grouper

import (
	"context"
	"errors"
	"fmt"
)

// ErrDivideByZero represents a situation where the passed divisor is zero in a division operation.
var ErrDivideByZero = errors.New("cannot divide by zero")

// Calculator is a driven port providing math as a service
type Calculator interface {
	// Divide divides the dividend by the divisor and returns the result.
	// Returns ErrDivideByZero if the divisor is zero.
	Divide(ctx context.Context, dividend, divisor int) (int, error)

	// ...other available basic math functions
}

// ErrBadNumberOfGroups reports an invalid number of groups
var ErrBadNumberOfGroups = errors.New("the number of groups must be positive and not zero!")

// ErrBadTotal reports an invalid total number of things to group
var ErrBadTotal = errors.New("must have 1 or more things")

// Core is a driving port for the grouper business logic
type Core interface {
	// GroupSize returns the size of a single group if totalThings is divided into numGroups. Returns
	// ErrBadNumberOfGroups if the number of groups is not 1 or more or ErrBadTotal if totalThings is not 1 or more
	GroupSize(ctx context.Context, totalThings, numGroups int, calc Calculator) (int, error)
}

type CoreLogic struct{}

// GroupSize implements Core for CoreLogic
func (CoreLogic) GroupSize(ctx context.Context, totalThings, numGroups int, calc Calculator) (int, error) {
	groupSize, calcErr := calc.Divide(ctx, totalThings, numGroups)
	if calcErr != nil {
		// errors.Is() can be used to differentiate errors, even if they're wrapped via fmt.Errorf()
		if errors.Is(calcErr, ErrDivideByZero) {
			return 0, ErrBadNumberOfGroups
		}

		// fmt.Errorf can be used to wrap an error with additional context
		return 0, fmt.Errorf("unexpected calculation problem occurred during division: %w", calcErr)
	}

	// ...rest of the business logic
}
```
</details>

## Writing REST controllers (driving adapters)

REST controllers are driving adapters that connect to the business logic's driving ports. They implement the `router.Controller`
interface (see [repository layout](Navigation and Repository Layout.md) for where that package is) so they can be
attached to a `router.Router` instance. They invoke business logic and convert DTO (Data Transfer Object) types to domain types defined in the
business logic by the business logic's driven port definition.

### Defining a route

Routes can be added to existing controllers by creating member functions on the controller which implement the [echo.HandlerFunc](https://pkg.go.dev/github.com/labstack/echo/v4#HandlerFunc)
interface. Those handler functions can then be added to routes via the `AttachRoutes()` function on the controller, which
exists for the controller to implement the `router.Controller` interface. Additionally, if you don't want to write the code
to unmarshal and validate your request body, you can wrap your route handlers in `router.AutoBindAndValidate()` Here's an example:

```go
package samplecontroller

import (
	"github.com/labstack/echo/v4"
	"example.com/sample/commonlib/router"
)

// SampleController is an example REST controller.
type SampleController struct {
	core sample.Core
}

// ...samplecontroller constructors and stuff

// AttachRoutes implements router.Controller for SampleController
func (ctrl SampleController) AttachRoutes(rtr *echo.Echo) {
	// This registers ExampleRoute() to execute on the /api/v1/sample endpoint
	rtr.POST("/api/v1/sample", ctrl.ExampleRoute)

	// Wrapping a route in router.AutoBindAndValidate() will take care of unmarshalling and validating
	// the request body for you. Validation only occurs if the body implements validation.Validatable.
	rtr.POST("/api/v1/auto-bind-sample", router.AutoBindAndValidate(ctrl.AutoBindExampleRoute))
}

// ExampleRoute contains logic for a route
func (ctrl SampleController) ExampleRoute(ctx echo.Context) error {
	// Add route logic here!
}

// AutoBindExampleRoute defines a route which doesn't need to manually deserialize the request body
func (ctrl SampleController) AutoBindExampleRoute(ctx echo.Context, requestBody SampleRequestBody) error {
	// Add route logic here! The request body is available via the requestBody parameter
}
```

### DTOs, validation, and responding

DTOs are defined in the same package as REST controllers. They exist so that any serialization specifics are abstracted
away from the business logic and so that fields from domain data can be excluded from the outside world if need be.

Requests should end with "Request" and responses should end with "Response" so you can tell them apart at a glance.
Responses should define JSON aliases so the resulting fields are encoded in **camelCase**. Here's an example:

```go
// SuperCoolActionRequest requests a super cool action on the REST controller
type SuperCoolActionRequest struct {
	// Note that SomeProperty doesn't need JSON field tags
	SomeProperty string
}

type SuperCoolActionResponse struct {
	// Result uses the "json" tag to make Result lowercase when it's serialized
	Result string `json:"result"`
}
```

Validation on requests can be done by implementing the [validation.Validatable](https://pkg.go.dev/github.com/jellydator/validation#Validatable)
interface provided by [jellydator/validation](https://github.com/jellydator/validation), a maintained fork of [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation).
Implementation of this interface enables nested struct validation and auto-validation via `router.AutoBindAndValidate()`,
so it's a good practice to implement it rather than using a custom function. Here's an example of adding validation on
the request from the example above:

```go
type SuperCoolActionRequest struct {
	SomeProperty string
}

// Validate implements validation.Validatable for SuperCoolActionRequest.
func (req SuperCoolActionRequest) Validate() error {
	return validation.ValidateStruct(&req,
		// This field validation states SomeProperty should be required and must have a length between 5 and 50 characters.
		validation.Field(&req.SomeProperty, validation.Required, validation.Length(5, 50)),
    )
}
```

**NOTE THAT YOU SHOULD USE [validation.Required](https://github.com/jellydator/validation#built-in-validation-rules) WHEREVER POSSIBLE.**
Go's JSON deserialization doesn't return errors on missing fields and fills them with their zero values, making
it very easy to get garbage data when you don't want it.

#### Omitting values in JSON output

The common library provides the `Nullable[T]` data structure to represent values which may or may not be present, and it
is compatible with JSON serialization and deserialization. The `IsPresent` field represents whether the value is
present or not, so be sure to check that before using the value in the `Value` field. Here's an example of usage:

```go
package sample

import example.com/sample/commonlib/types"

type NameWithOptionalLastName struct {
    FirstName string `json:"firstName"`
    LastName types.Nullable[string] `json:"lastName"`
}
```

An associated validation rule, `types.NullablePresent` is available for conditionally requiring that a nullable value is
present based on other factors, such as the values of other fields in the data structure. It can be used like this:

```go
package sample

import (
	"github.com/jellydator/validation"
	example.com/sample/commonlib/types"
)

// ...continuing from the previous example

func (name NameWithOptionalLastName) Validate() error {
	// Last name will be required if the first name is longer than 10 characters
	var lastNameValidations []validation.Rule

	if len(name.FirstName) > 10 {
        // Note that the validation rule's generic type must match the underlying type in the nullable
		lastNameValidations = append(lastNameValidations, types.NullablePresent[string]{})
	}
	
    return validation.ValidateStruct(&name,
        validation.Field(&name.LastName, lastNameValidations...),
    )
}
```

#### Automatically binding and validating

To avoid boilerplate, the function `router.AutoBindAndValidate()` can be used to automatically unmarshal the request body
and validate the contents, if applicable. [Methods on the echo.Context](https://echo.labstack.com/docs/response) such as
`Context.JSON()` and `Context.NoContent()` can be used to respond. Here's an example:

```go
package samplecontroller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"example.com/sample/commonlib/router"
)

type SampleController struct {
	// ...controller contents
}

func (ctrl SampleController) AttachRoutes(rtr *echo.Echo) {
	// AutoBindAndValidate will automatically unmarshal the request body for you in the second parameter of
	// the passed function. If the request body also implements validation.Validatable, it will be validated.
	rtr.POST("/api/v1/auto-bind-sample", router.AutoBindAndValidate(ctrl.ExampleAutoBindRoute))
}

func (ctrl SampleController) ExampleAutoBindRoute(ctx echo.Context, actionRequest SuperCoolActionRequest) error {
	// ...call into business logic here

	// Pretend "actionResponse" is produced after invoking the business logic
	return ctx.JSON(http.StatusOK, SuperCoolActionResponse {
		Response: actionResponse,
    })
}
```

#### Manually binding and validating

If you want to perform data binding and validation manually, you can use `request.Bind` to deserialize the body of the
request and invoke `.Validate()` on the request to make sure it's valid. [Methods on the echo.Context](https://echo.labstack.com/docs/response) such as
`Context.JSON()` and `Context.NoContent()` can be used to respond. Here's an example:

<details>
<summary>Request and response example</summary>

```go
package samplecontroller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"example.com/sample/commonlib/request"
)

type SampleController struct {
	// ...controller contents
}

// ...AttachRoutes() implementation

func (ctrl SampleController) ExampleRoute(ctx echo.Context) error {
	// Define the request with the zero value
	var requestBody SuperCoolActionRequest

	// Unmarshal the request body with request.Bind, passing a pointer to the request body so Bind can fill its contents
	if unmarshalErr := request.Bind.BindBody(ctx, &requestBody); unmarshalErr != nil {
		// Handle the deserialization error (see "canned error responses")
	}

	// Invoke the validation function on the request to make sure its contents are valid
	if validationErr := requestBody.Validate(); validationErr != nil {
		// Handle the validation error (see "canned error responses")
	}

	// ...rest of the controller

	// "actionResponse" is produced after invoking business logic
	return ctx.JSON(http.StatusOK, SuperCoolActionResponse {
		Response: actionResponse,
	})
}
```
</details>

#### Returning non-nil slices

If you declare a slice using the zero value, the slice will be set to `null` by default. To avoid this, use the `response.NonNilSlice()` function
from the common library to ensure your slices actually show up as empty slices and not null to API consumers.

```go
package sample

import "example.com/sample/commonlib/response"

type MyCoolResponse struct {
	StringSlice []string `json:"stringSlice" example:"[\"abc\", \"def\"]"`
}

// Assume this function generates a DTO from a domain object
func MyResponseFromDomain(stringSlice []string) MyCoolResponse {
	return MyCoolResponse{
		StringSlice: response.NonNilSlice(stringSlice),
	}
}
```

### Canned error responses

In order to make error responses consistent across the application and reduce the amount of boilerplate necessary, several
error response functions are available to quickly respond with different HTTP status codes and report an error. These
error response functions are available in the `response` package (see [repository layout](Navigation and Repository Layout.md)
for where that is).

The `response.APIErrorHelper` type returned from the canned response functions can be customized after its construction
so that you can customize the error message in the standard error response type.

Here's a usage example, responding with a 400 bad request when a request is invalid:

```go
// ...controller implementation code

func (ctrl SampleController) ExampleRoute(ctx echo.Context) error {
	var requestBody SuperCoolActionRequest

	// ...deserialize the request

	if validationErr := requestBody.Validate(); validationErr != nil {
		// Never a bad idea to log a warning on an expected error (see the logging docs for more info on log levels)
		logger.Log.Warn("Super cool action request was invalid.", zap.Error(validationErr))
		// Create the "Bad Request" canned response
		cannedResponse := response.BadRequest(validationErr)
		// Customize the error message by editing the "Description" field, though this is optional
		cannedResponse.Description = "Your super cool action request was invalid."
		// Respond by calling the Respond() function
		return cannedResponse.Respond(ctx)
    }
}
```

### Driven port management

Driven ports should be accepted by the REST controller in its constructor so that they can be passed along to Driving ports.
The controller can add the driving port implementation automatically and provide a separate, private constructor to swap
out the driving port for a mock in tests. More information on that private constructor can be found in the [testing docs](Testing.md#testing-rest-controllers).

Here's an example:

```go
type SampleController struct {
	// Driven port which will be passed to grouper.Core.GroupSize (see "business logic" example for reference)
	calculator grouper.Calculator

	// Driving port to invoke business logic with
	core grouper.Core
}

// Constructor for the SampleController
func New(calculator sample.Calculator) SampleController {
	return SampleController {
		calculator: calculator,
		core: sample.CoreLogic{},
    }
}

func newWithCore(core grouper.Core) SampleController {
	return SampleController {
		core: core,
		// Other fields can be omitted because the driven port is mocked
    }
}
```

### Initiating database transactions across business logic

Since business logic must have implementation details abstracted away, business-logic-wide database transactions must
be triggered in driving ports using the `database.WithTransaction()` or `database.WithTransactionReturning()` function. It alters the current request's `context.Context`
carrying the database connection to initiate a transaction, passing the modified context to the provided function. Here's an example,
re-using our `grouper` example from the [business logic](#errors-in-business-logic) section:

```go
// ctrl.Core is our instance of grouper.Core and ctx is our echo.Context instance
// We need to extract the context.Context from the echo.Context first
requestCtx := request.ExtractContext(ctx)
// database.WithTransactionReturning will return the return value of the passed function
computationResult, computationErr := database.WithTransactionReturning(requestCtx, func (transactionCtx context.Context) (int, error) {
	return ctrl.Core.GroupSize(transactionCtx, 15, 5)
})
```

If the function passed to `database.WithTransaction()` or `database.WithTransactionReturning()` returns an error, the database transaction
wrapping any inserts triggered by driven ports will be automatically rolled back. `database.WithTransactionReturning()` slightly differs from
`database.WithTransaction()` because it allows one to return a return value from the passed function, which will then be returned by `database.WithTransactionReturning()`.

### Attaching controllers to the router

REST controllers implementing the `router.Controller` interface can be attached to the `router.Router` instance via the
`router.Router.AttachControllers()` function. Each controller is instantiated in its own function in `bootstrap.go` (see [repo layout docs](Navigation and Repository Layout.md))
and added to the slice of controllers, which is passed along to the `AttachControllers()` function. Here's a small example,
using the sample controller as seen in the previous examples:

```go
package main

import "example.com/sample/commonlib/router"

// This instantiates the router
func CreateRouter() router.Router {
	router := router.New()
	controllers := CreateControllers()
	router.AttachControllers(controllers)

	return router
}

// This instantiates the full list of controllers
func CreateControllers() []router.Controller {
	controllers := []router.Controller{
		sample(),
	}

	return controllers
}

// This instantiates the sample controller, initializing any driven ports and passing them to the controller constructor
func sample() sample.SampleController {
	// This is the implementation of the grouper.Calculator interface seen in the "business logic" section
	grouper := adapter.CalculatorAsAService{}
	return sample.New(grouper)
}
```

## Connecting to external data sources (driven adapters)

Driven adapters are called by the business logic to reach external systems. These adapters may connect to other microservices,
the database, message queues, or other similar things.

### Acquiring a database connection

When writing a driven adapter that connects to the database, you can use the database connection in the current `context.Context`
via the extractor function `database.RetrieveFromContext()`. This will return a `database.Connection`, which is a generalized
interface over both [sqlx.DB](https://pkg.go.dev/github.com/jmoiron/sqlx#DB) and [sqlx.Tx](https://pkg.go.dev/github.com/jmoiron/sqlx#Tx)
so you can write your database adapter with a single implementation regardless of whether the current operation is in a transaction or not.

Here's an example of extracting the connection:

```go
type DatabaseAdapter struct {}

func (DatabaseAdapter) GetSomeData(ctx context.Context) (int, error) {
    connection := database.RetrieveFromContext(ctx)

    // Now you can do any database queries you want with "connection"
}
```

The database connection is attached to the context via the database context middleware, as described in [the middleware docs](Middleware.md#database-connection-middleware).
You can see all the querying options available on the `database.Connection` type (see the [repo layout](Navigation and Repository Layout.md)
for where the database package is).

### Database-specific DTOs

It is highly recommended to extract database query results into database-specific DTOs so database types in the data structure
such as [sql.NullString](https://pkg.go.dev/database/sql#NullString) don't leak into the business logic. You'll just need
to convert these DTOs into domain types before they're returned to the business logic.

Note also that you should alias the columns you extract from the database, as the Go SQL library only works off of columns
in the results and not table names as described in [this section of the SQLX guide](https://jmoiron.github.io/sqlx/#altScanning).
If you need columns to be a different name from the aliases in the database query, you can use the `db` tag on your struct
fields, similarly to using the `json` tag in REST controllers.

Here's an example:

```go
type DatabaseAdapter struct {}

type ExtractedDataDto struct {
	SomeStringResult sql.NullString `db:"stringResult"`
	SomeNumberResult int            `db:"numberResult"`
}

func (DatabaseAdapter) GetSomeData(ctx context.Context) (int, error) {
	connection := database.RetrieveFromContext(ctx)

	// Create your DTO with the zero value
	var dbDto ExtractedDataDto
	// See the documentation comment on database.Connection.Get for why we're using it here
	dbErr := connection.Get(&dbDto, `
        select sampleTable.stringValue as stringResult, sampleTable.intValue as numberResult
            from sampleTable
            limit 1
    `)
	// error handling...

	return dbDto.SomeNumberResult
}
```

### Communicating with other systems over HTTP

TBD, we can take care of this subsystem in another ticket. Needs to be done in a way that we can mock responses from external systems.