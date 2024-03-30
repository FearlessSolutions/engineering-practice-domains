package request

import (
	"github.com/labstack/echo/v4"
)

// Bind is a globally accessible instance of an echo.DefaultBinder, making it easy to
// specifically parse certain parts of an incoming HTTP request into different data structures.
var Bind = new(echo.DefaultBinder)
