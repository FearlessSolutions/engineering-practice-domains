package middleware

import (
	"example.com/sample/commonlib/config"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// StandardMiddleware returns a set of pre-configured middleware that can be applied across many microservices
func StandardMiddleware(options config.Registry, databaseConnection *sqlx.DB) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Recover(),
		CorsMiddleware(options),
		LoggingMiddleware(),
		AuthMiddleware(),
		DatabaseContextMiddleware(databaseConnection),
	}
}
