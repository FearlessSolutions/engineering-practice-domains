package request

import (
	"context"

	"github.com/labstack/echo/v4"
)

// ExtractContext retrieves a context.Context from an echo.Context
func ExtractContext(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}
