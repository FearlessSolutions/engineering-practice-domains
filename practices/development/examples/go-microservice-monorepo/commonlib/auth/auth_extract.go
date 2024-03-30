package auth

import (
	"github.com/labstack/echo/v4"
)

func RetrieveAuthClaims(c echo.Context) (CustomClaims, bool) {
	rawClaims := c.Get("user")

	if rawClaims != nil {
		claims := rawClaims.(CustomClaims)
		return claims, true
	} else {
		return CustomClaims{}, false
	}
}
