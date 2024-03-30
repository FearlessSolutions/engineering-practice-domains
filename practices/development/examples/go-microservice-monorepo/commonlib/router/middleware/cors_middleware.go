package middleware

import (
	"strings"

	"example.com/sample/commonlib/config"
	"example.com/sample/commonlib/config/sharedoptions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CorsMiddleware is a middleware that auto-handles the CORS preflight a browser sends to only allow
// certain frontend origins to access our API. It pulls the list of allowed origins from sharedoptions.AllowedOrigins
// but if that's not present it uses localhost:3000 as the only acceptable origin
func CorsMiddleware(options config.Registry) echo.MiddlewareFunc {
	var allowedOrigins []string
	rawOriginsList, listPresent := options.Get(sharedoptions.AllowedOrigins)
	if listPresent {
		allowedOrigins = strings.Split(rawOriginsList, ",")
	} else {
		allowedOrigins = []string{"http://localhost:8080"}
	}

	// TODO tweak the CORS headers as necessary.
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     middleware.DefaultCORSConfig.AllowMethods,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	})
}
