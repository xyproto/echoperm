package echoperm

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
)

// Middleware sets up a middleware handler for Echo.
//
// Requires a Redis hostname (can be "localhost"), a Redis password (can be
// blank) and a custom message for when permissions are denied.
func Middleware(redisHostname, redisPassword, denyMessage string) (echo.MiddlewareFunc, pinterface.IUserState, error) {
	userstate, err := permissions.NewUserStateWithPassword2(redisHostname, redisPassword)
	if err != nil {
		return nil, nil, err
	}
	perm := permissions.NewPermissions(userstate)

	return echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			// Check if the user has the right admin/user rights
			if perm.Rejected(c.Response(), c.Request()) {
				// Deny the request
				return echo.NewHTTPError(http.StatusForbidden, denyMessage)
			}
			// Continue the chain of middleware
			return next(c)
		})
	}), userstate, nil
}
