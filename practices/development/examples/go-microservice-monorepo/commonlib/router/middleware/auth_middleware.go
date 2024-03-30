package middleware

import (
	"example.com/sample/commonlib/auth"
	"example.com/sample/commonlib/response"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware constructs a middleware function that just parses the claims out of a JWT and doesn't verify them, as
// an Istio proxy verifies auth tokens before requests hit this microservice in production. The auth info this middleware
// attaches to the request can be easily extracted from the echo context with auth.RetrieveAuthClaims
func AuthMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		// Skipper determines if the app should skip checking/parsing the JWT
		Skipper: func(ctx echo.Context) bool {
			// TODO add auth skip logic, if applicable
			// Add code here to conditionally skip JWT validation.
			// Just so you can play with the app locally, all JWT validation is skipped here.
			return true
		},

		// ErrorHandler allows us to return a custom error when the JWT is bad
		ErrorHandler: func(c echo.Context, err error) error {
			return response.UnauthorizedWithErr(err).Respond(c)
		},

		// TODO if your app needs to handle auth in a special way, use ParseTokenFunc. This function is currently
		//   configured to parse a JWT using the auth.RawJWTCustomClaims struct to capture relevant JWT payload fields
		//   and to not verify the JWT's signature. This should be changed depending on your app's security requirements.
		ParseTokenFunc: func(c echo.Context, authString string) (interface{}, error) {
			var claims auth.RawJWTCustomClaims
			parser := jwt.NewParser()
			if _, _, parseErr := parser.ParseUnverified(authString, &claims); parseErr != nil {
				return nil, parseErr
			}

			return claims.ToCustomClaims(), nil
		},

		// TODO uncomment this or use the "SigningKeys" field for JWT verification
		//SigningKey: []byte(""),
	}

	return echojwt.WithConfig(config)
}
