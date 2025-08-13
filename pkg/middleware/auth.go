package middleware

import (
	"monitoring-service/pkg/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
			}

			claims, err := utils.ParseJWTToken(tokenString, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Set user info ke context
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role").(string)

			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}
	}
}
