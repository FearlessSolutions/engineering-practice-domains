# Authentication

In production, JWT validation and verification is handled by the Istio proxy between all services in the
Kubernetes cluster. For this reason, the [authentication middleware](Middleware.md#auth-middleware)
simply deserializes the JWT claims if a JWT is passed in the `Authorization` header of an incoming HTTP
request, but doesn't otherwise enforce use of credentials when hitting the API.

## Extracting authentication information in a REST controller

The `auth` package provides a JWT extractor function to retrieve the claims parsed by the authentication middleware.
This function is called `auth.RetrieveKeycloakClaims()`, and because the authentication token may not always be present,
it returns a boolean in addition to the returned data structure stating whether the token is present or not. It is recommended
that you respond with a 401 Unauthorized if you expect the token to be present. See the [microservice architecture docs](Microservice Architecture.md#writing-rest-controllers-driving-adapters)
for information on writing REST controllers.

Here's an example of using `auth.RetrieveKeycloakClaims()`:

```go
package sample

import (
	"github.com/labstack/echo/v4"
	"example.com/sample/commonlib/auth"
	"example.com/sample/commonlib/response"
)

type SampleController struct {
	// ...controller members
}

// ...SampleController constructor and AttachRoutes implementation

func (ctrl SampleController) ExtractsTokenData(ctx echo.Context) error {
	claims, claimsArePresent := auth.RetrieveAuthClaims(ctx)
	if !claimsArePresent {
		return response.Unauthorized().Respond(ctx)
	}

	// ...rest of the route implementation
}
```
