package delivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/bekaza/go-clean/domain"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware ...
type AuthMiddleware struct {
	userService domain.UserService
}

// NewAuthMiddleware ...
func NewAuthMiddleware(us domain.UserService) *AuthMiddleware {
	return &AuthMiddleware{us}
}

// CORS will handle the CORS middleware
func (am *AuthMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// AuthRequire ...
func (am *AuthMiddleware) AuthRequire(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Auth require",
			})
			return nil
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Auth is invalid",
			})
			return nil
		}

		if headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Auth is invalid",
			})
			return nil
		}

		user, err := am.userService.ParseToken(c.Request().Context(), headerParts[1])
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "invalid access token" {
				status = http.StatusUnauthorized
			}

			c.JSON(status, echo.Map{
				"message": "invalid access token",
			})
			return nil
		}
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, domain.CtxUserKey, user)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
