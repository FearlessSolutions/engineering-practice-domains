# Testing

We use the [testify](https://github.com/stretchr/testify) library to improve the tools we have available to write tests.
Testify provides the capability to write test suites, perform more powerful assertions, and even perform parametrized
tests in a really easy way.

Test files should be defined next to the file they intend to test, and are detected by the Go testing tool via the [_test suffix](https://go.dev/doc/tutorial/add-a-test)
in the file name. Tests can then be run with this command from the repository root:

```bash
go test ./...
```

## Test suites

Generally, you'll want to use test suites over the normal Go test functions to write your tests. Suites allow you to
write setup and teardown code which run before and after sets of tests so that you don't need to write lengthy setup
code or cleanup logic for every test. A suite is defined by embedding a [suite.Suite](https://pkg.go.dev/github.com/stretchr/testify/suite#Suite)
into a struct in your test file. You then need to define a traditional Go test and tell it to run the suite. Test functions
on the suite are then detected as member functions of the suite with the prefix `Test`.

Here's an example of a full test suite setup:

```go
type NumbersMakeSenseSuite struct {
	// Embed the suite.Suite in this struct
	suite.Suite
}

// Create the traditional test function that kicks off the suite
func TestNumbersMakeSenseSuite(t *testing.T) {
	suite.Run(t, new(NumbersMakeSenseSuite))
}

// Now, create a test (a member function of the suite prefixed by "test")
func (suite *NumbersMakeSenseSuite) TestNumbersEqualThemselves() {
	suite.Require().Equal(5, 5)
}
```

### Asserting

Embedding the `suite.Suite` inside your test suite adds all the member functions of `suite.Suite` to your test suite,
including assertion functions.

There are two assertion functions on a `suite.Suite`, and they behave slightly differently:

1. `suite.Suite.Assert()` - Assertions made via this function will allow the test to keep running but fail the test at the end - this is called a "soft assertion"
    * These assertions return a boolean value which is true if the assertion passed
2. `suite.Suite.Require()` - Assertions made via this function will immediately stop the test if they fail. **This is the preferred method, unless you really need to perform soft assertions.**

All the assertions defined in `testify`'s [require](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/require) and [assert](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/assert)
assertions are available off of the respective suite functions.

### Setup and teardown functions

Test suites can contain additional pieces of data accessible from each test, and that data can be set up via setup and
teardown functions which run before or after different points of tests.
This is done by implementing certain interfaces on the `suite.Suite`:

1. [SetupAllSuite](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#SetupAllSuite) defines a function which will be executed at the beginning of the entire suite. This is a good place to initialize the global logger (see [the logging docs](Logging.md#instantiating-the-logger-directly) for more info)
2. [SetupTestSuite](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#SetupTestSuite) defines a function which will be executed right before every test. This is a good place to initialize mocks, but not necessarily define expectations.
3. [SetupSubTest](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#SetupSubTest) defines a function which runs before each subtest in a suite, which is useful if you have [parametrized tests](#parametrized-tests).
4. [TearDownAllSuite](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#TearDownAllSuite) is just like `SetupAllSuite` except its function runs after the suite completes
5. [TearDownTestSuite](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#TearDownTestSuite) is just like `SetupTestSuite` except its function runs after every test, which is a good place to finalize mock controllers.
6. [TearDownSubTest](https://pkg.go.dev/github.com/stretchr/testify@v1.8.4/suite#TearDownSubTest) is just like `SetupSubTest` but it runs after each subtest

Here's an example of implementing the `SetupTestSuite` interface to do something before every test:

```go
type NumbersMakeSenseSuite struct {
	suite.Suite
	// Pretend this is a mock or something
	numberToCompareAgainst int
}

// ...suite run function

// SetupTest implements suite.SetupTestSuite for NumbersMakeSenseSuite
// This function will now be run before every test
func (suite *NumbersMakeSenseSuite) SetupTest() {
	fmt.Println("Numbers definitely make sense!!!")
	suite.numberToCompareAgainst = 5
}

func (suite *NumbersMakeSenseSuite) TestNumbersEqualThemselves() {
	suite.Require().Equal(suite.numberToCompareAgainst, 5)
}
```

### Parametrized tests

Especially for validation, sometimes you need to run the same test logic with a bunch of different pieces of data.
Instead of writing a test over and over again, you can instead write parametrized tests which just run the same test
repeatedly but with different data every time. Parametrized tests can be run in a loop using `suite.Suite.Run()`, here's
an example:

```go
func (suite *NumbersMakeSenseSuite) TestNumbersLessThanFive() {
	// You can actually define types inside a function, which is ideal for packaging test parameters:
	type numberCompareCase struct {
		testName string
		value int
    }

	// Now, define the parameters for the test cases:
	testCases := []numberCompareCase{
		{
			testName: "4 is less than 5",
			value: 4,
        },
		{
			testName: "3 is less than 5",
			value: 3,
        },
		{
			testName: "2 is less than 5",
			value: 2,
        },
		{
			testName: "1 is less than 5",
			value: 1,
        }
    }

	for _, testCase := range testCases {
		suite.Run(testCase.testName, func() {
			suite.Require().Less(testCase.value, 5)
        })
    }
}
```

## Mocking and Faking

Because Go is a compiled language, we can't swap out real struct types at runtime like we can with languages like Java
or TypeScript. Instead, we need to create mocks which implement interfaces we use in the code, such as the ports we define
between the business logic and other layers of the application.

One strategy we could employ is creating **Fakes**, which would be manual implementations of interfaces which mimic the
behavior of an adapter a piece of code depends on. This can be useful for cases where business logic performs a ton of
operations which change the behavior of the adapter based on the data passed to it. It does take development effort to
write, though.

If we just need a quick way to get a function to respond a certain way, we can use mocks generated by Uber's
**MockGen** tool supplied with their [gomock](https://github.com/uber-go/mock) library.

### Using mocks

Existing mocks can be found in MockGen-generated files with the `_mocks` suffix in their filenames. See the next section
on how these are generated.

All mocks are instantiated by passing an instance of [gomock.Controller](https://pkg.go.dev/go.uber.org/mock@v0.3.0/gomock#Controller)
which is used to verify expected numbers of calls, hold a reference to the [testing.T](https://pkg.go.dev/testing#T) instance,
and for "call in order" verification. The controller should be set up before each test, and its "Finish" method should be
invoked after each test.

Here's an example with an imaginary "SuperCoolAdapterMock":

```go
package sample

import (
   "github.com/stretchr/testify/suite"
   "go.uber.org/mock/gomock"
)

type SuperCoolBusinessLogicSuite struct {
   suite.Suite
   mockController *gomock.Controller
   adapterMock *SuperCoolAdapterMock
}

// ...test run function

// SetupTest runs before every test and sets up the mock and mock controller
// Expectations should be set in the test, but mocks can be set up here
func (suite *SuperCoolBusinessLogicSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.adapterMock = SuperCoolAdapterMock(suite.mockController)
}

// TearDownTest runs after every test and verifies expected functions in the test actually got called
func (suite *SuperCoolBusinessLogicSuite) TearDownTest() {
	suite.mockController.Finish()
}
```

Expectations and mocked return values can then be expressed with a mock's `EXPECT()` function:

```go
func (suite *SuperCoolBusinessLogicSuite) TestSuperCoolOperationPerformsSuccessfully() {
	suite.adapterMock.EXPECT().
		// This says we expect SomeMemberFunction() on the adapter to be called
		// The "gomock" package has argument matchers, or you can just pass values to be compared against
		SomeMemberFunction(gomock.Any(), 12345).
		// The Return function specifies the return value on function invocation
		Return(54321, nil)

	returnValue, someError := SomeBusinessLogic{}.DoOperation(context.Background(), 12345, suite.adapterMock)
	// ...rest of the test
}
```

Since `EXPECT()` returns a `gomock.Call` instance, you can see other things you can do with a mock on [the relevant documentation](https://pkg.go.dev/go.uber.org/mock@v0.3.0/gomock#Call).

### Generating mocks

We can generate mocks using the `mockgen` tool provided with Uber's mock library, which integrates nicely with the `go generate`
command. In order to generate mocks, you'll need to first install the `mockgen` tool. You can do so with the following
command:

```bash
go install go.uber.org/mock/mockgen@latest
```

Any existing mocks can then be regenerated by running the following command at the repository root:

```bash
go generate ./...
```

Note that that command will fail if you've made manual edits to the implementation of generated mocks since the last time.

"Go generate comments" in the code tell `go generate` how to run `mockgen` to generate our mocks. These comments will
have the following format:

```go
//go:generate mockgen -destination ./FILENAME_mocks.go -package PACKAGE_NAME . INTERFACES,TO,IMPLEMENT,COMMA,SEPARATED
```

For example, in the file "sample.go" in the "coolstuff" package generating mocks for two interfaces,
the code would look like this:

```go
package coolstuff

// The lack of a space at the beginning of the comment is important!!

//go:generate mockgen -destination ./sample_mocks.go -package coolstuff . MyInterfaceToBeMocked,AnotherInterfaceToBeMocked

type MyInterfaceToBeMocked interface {
	DoSomethingCool() (string, error)
}

type AnotherInterfaceToBeMocked interface {
	DoSemethingRadical() (int, error)
}
```

Now that that comment is added, `go generate ./...` at the root of the repository will generate a `sample_mocks.go` file
containing implementations for `MyInterfaceToBeMockedMock` and `AnotherInterfaceToBeMockedMock`.

## Testing REST controllers

Testing REST controllers is a little tricky because you need an `echo.Context` to execute REST controller functions.
Thankfully, there's a `testhelper` package in the `request` package of `common_lib` which contains a test request builder
specifically for creating `echo.Context` instances and reading HTTP responses from said controller function.

You can instantiate a `testhelper.RequestBuilder` via `testhelper.NewRequestBuilder()`, which initially just accepts the
HTTP verb and the URI the request is being made to. Other functions on the builder will allow you to add authentication
information, HTTP request bodies, and inject request path parameters too.

For REST controllers that require authentication, you can construct a mock token via `auth.MockKeycloakClaims()`. The mocked struct should be passed to `testhelper.WithAuth()` when constructing a mock request.

Remember that in [the architecture docs](Microservice Architecture.md#driven-port-management) it's mentioned that you
should have a private constructor for testing that accepts a mock for the business logic. We'll use that here for a
controller that allows you to look up a greeting from a collection of greetings - `GET /api/v1/sample/greetings/:id`

Note also that the `response` package also has a `testhelper` subpackage which helps with extracting the body of a
controller's response via `testhelper.UnmarshalBody()` (preventing you from needing to do error handling while unmarshaling).

Here's how the test helpers might look together (in this example, we alias the two imports as `reqhelper` and `reshelper`):

<details>
<summary>REST Controller test example code</summary>

```go
package sample

import (
   "github.com/labstack/echo/v4"
   "github.com/stretchr/testify/suite"
   "go.uber.org/mock/gomock"
   "net/http"
   "example.com/sample/commonlib/auth"
   reqhelper "example.com/sample/commonlib/request/testhelper"
   reshelper "example.com/sample/commonlib/response/testhelper"
)

type SampleControllerSuite struct {
   suite.Suite
   mockController    *gomock.Controller
   mockBusinessLogic *MockCore
}

// ...setup & teardown functions, suite start function

func (suite *SampleControllerSuite) TestSampleControllerSuccess() {
   // Set up mock returns for any business logic the controller invokes
   suite.mockBusinessLogic.EXPECT().RetrieveGreeting().Return("Hello", nil)

   // request is an echo.Context, responseRecorder records the HTTP response and lets us extract it later
   request, responseRecorder, buildErr := reqhelper.NewRequest(echo.GET, "/api/v1/sample/greetings/:id").
	    // WithAuth adds the specified authentication information to the request being built
	   WithAuth(auth.MockKeycloakClaims()).
        // WithPathParams mimics the extraction echo would have done for us
        // Other "With*" functions are available on the builder for attaching info to the request too
	   WithPathParams(map[string]string{
              "id": "1",
           }).
      // Build will construct everything we need to make the request
      Build()
   // Make sure the request was constructed successfully
   suite.Require().NoError(buildErr)

   // Now we can construct the controller with the private constructor
   ctrl := newWithCore(suite.mockBusinessLogic)
   resErr := ctrl.GetGreetingById(request)
   // Although echo.HandlerFunc says it returns an error, we really shouldn't be returning errors
   suite.Require().NoError(resErr)

   // Now we can extract the response from the response recorder and inspect it
   response := responseRecorder.Result()
   // Make sure we got the right status code
   suite.Require().Equal(http.StatusOK, response.StatusCode())
   // Use the response helper to extract the response
   var responseBody GreetingByIdResponse
   // We pass the suite here because UnmarshalBody will auto-fail the test on an unparsable body
   reshelper.UnmarshalBody(&suite.Suite, response.Body, &responseBody)
   // Now we can inspect the body returned in the response
   suite.Require().Equal("Hello", responseBody.Greeting)
}
```

</details>

And here's an example using an [auto-bind route](Microservice Architecture.md#automatically-binding-and-validating):

<details>
<summary>Auto-bind route testing example</summary>

```go
package sample

import (
   "github.com/labstack/echo/v4"
   "github.com/stretchr/testify/suite"
   "go.uber.org/mock/gomock"
   "net/http"
   reqhelper "example.com/sample/commonlib/request/testhelper"
   reshelper "example.com/sample/commonlib/response/testhelper"
   "example.com/sample/commonlib/router"
)

type SampleControllerSuite struct {
   suite.Suite
   mockController    *gomock.Controller
   mockBusinessLogic *MockCore
}

// ...setup & teardown functions, suite start function

func (suite *SampleControllerSuite) TestSampleControllerSuccess() {
   // Set up mock returns for any business logic the controller invokes
   suite.mockBusinessLogic.EXPECT().
      AddGreeting(gomock.Any(), gomock.Any()).
      Return(nil)

   // Set up the request body
   reqBody := SampleGreetingRequest{
      Name: "Jason",
   }
   // request is an echo.Context, responseRecorder records the HTTP response and lets us extract it later
   request, responseRecorder, buildErr := reqhelper.NewRequest(echo.POST, "/api/v1/sample/greetings").
      // WithBody adds the passed data structure as the request body
      WithBody(reqBody).
      // Build will construct everything we need to make the request
      Build()
   // Make sure the request was constructed successfully
   suite.Require().NoError(buildErr)

   // Now we can construct the controller with the private constructor
   ctrl := newWithCore(suite.mockBusinessLogic)
   // We'll need to wrap the target route in AutoBindAndValidate
   route := router.AutoBindAndValidate(ctrl.AddGreeting)
   // Now we invoke the wrapped route
   resErr := route(request)
   // Although echo.HandlerFunc says it returns an error, we really shouldn't be returning errors
   suite.Require().NoError(resErr)

   // Now we can extract the response from the response recorder and inspect it
   response := responseRecorder.Result()
   // Make sure we got the right status code
   suite.Require().Equal(http.StatusCreated, response.StatusCode())
}
```
</details>


## Testing database-based driven adapters

Because the raw `sqlx.DB` or `sqlx.Tx` are masked by the `databese.Connection` interface when invoking `database.RetriveFromContext()`
(see [the microservice architecture docs](Microservice Architecture.md#acquiring-a-database-connection) for more info),
a mock of the `database.Connection` interface is exposed for testing via `database.MockConnection`. You'll just need to
inject it into the context passed to your database adapter implementation via `database.CreateDerivativeMockContext`.

Here's an example of how to use it all together:

<details>
<summary>Database adapter test example</summary>

```go
package sample

import (
   "context"
   "github.com/stretchr/testify/suite"
   "go.uber.org/mock/gomock"
   "example.com/sample/commonlib/database"
)

type SampleSuite struct {
   suite.Suite
   mockController *gomock.Controller
   mockConnection *database.MockConnection
   connectionCtx  context.Context
}

// ...suite runner function

func (suite *SampleSuite) SetupTest() {
   suite.mockController = gomock.NewController(suite.T())
   suite.mockConnection = database.NewMockConnection(suite.mockController)
   // This is where the mock connection gets injected into the context
   suite.connectionCtx = database.CreateDerivativeMockContext(context.Background(), suite.mockConnection)
}

func (suite *SampleSuite) TearDownTest() {
   suite.mockController.Finish()
}

func (suite *SampleSuite) TestLogicAroundDatabaseWorks() {
   // Set up mock return on database select. In the function it selects to a list of strings
   suite.mockConnection.EXPECT().
      Select(gomock.Any(), gomock.Any()).
      // The SetArg here is important! This allows your mock to write a value to the pointer passed as
      // the destination in the first parameter of database.Connection.Select()
      SetArg(0, []string{"A", "B", "C"}).
      // Select then returns an error, in this case we'll make it "nil"
      Return(nil)

   // Gotta make sure to pass the connection context defined in the test setup to get the "database connection" to the function
   listOfStrings, listErr := DatabaseAdapter{}.GetSomeData(suite.connectionCtx)
   suite.Require().NoError(listErr)
   suite.Require().ElementsMatch([]string{"A", "B", "C"}, listOfStrings)
}
```
</details>

## Testing functions requiring environment variables

*See the [configuration docs](Configuration.md#creating-a-configuration-registry-and-registering-options) for information
on configuration and configuration registries.*

To enforce a specific set of environment variables during a test, you can use the `config.NewMockRegistryBuilder()` function
which accepts a set of environment variables to use via a `map[string]string`. The mock registry builder will construct a real
`config.Registry`, but the constructed registry will use the map you pass rather than real environment variables. You'll
then need to assign your created registry to `options.Registry`.

It is recommended that you make use of `Option.VariableName()` when passing environment variables to registries in tests
as this will make the test resilient to renaming the actual environment variable.

Here's an example of a test that needs environment variables:

```go
func (suite *SampleSuite) TestWithRegistry() {
   registryBuilder := config.NewMockRegistryBuilder(map[string]string{
	   options.SomeConfigOption.VariableName(): "abcde",
	   sharedoptions.AnotherConfigOption.VariableName(): "defghi",
   })

   registryBuilder.AddOption(options.SomeConfigOption)
   registryBuilder.AddOption(sharedoptions.AnotherConfigOption)
   registry, buildErr := registryBuilder.VerifyAndBuild()
   suite.Require().NoError(buildErr)

   options.Registry = &registry

   // Now you can do your test
}
```
