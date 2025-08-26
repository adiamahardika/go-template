package middleware

import (
    "errors" 
    "monitoring-service/pkg/utils"
    "net/http"
    "strings"

    "github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int      `json:"user_id"`
	Roles  []string `json:"roles"`
}	

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

			var userRole string
			if len(claims.Roles) > 0 {
				userRole = claims.Roles[0]
			} else {
				userRole = ""
			}
			c.Set("role", userRole)

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

func GetUserIDFromToken(c echo.Context) (int, error) {
    user := c.Get("user_id")
    if user == nil {
        return 0, errors.New("user_id not found in context")
    }

    userID, ok := user.(int)
    if !ok {
        if userIDFloat, ok := user.(float64); ok {
            return int(userIDFloat), nil
        }
        return 0, errors.New("user_id is not of a valid type in context")
    }

    return userID, nil
}