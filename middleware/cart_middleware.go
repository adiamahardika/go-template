// middleware/jwt_roles.go
package middleware

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"monitoring-service/pkg/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

const ContextUserIDKey = "userID"

func JWTRequireRoles(cfg *config.Config, allowedRoles ...string) echo.MiddlewareFunc {
	roleAllowed := func(roles []string) bool {
		for _, r := range roles {
			for _, a := range allowedRoles {
				if r == a {
					return true
				}
			}
		}
		return false
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return helpers.StandardResponse(c, http.StatusUnauthorized, "Missing Authorization header", nil, nil)
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return helpers.StandardResponse(c, http.StatusUnauthorized, "Invalid Authorization header", nil, nil)
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWT.Secret), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					return helpers.StandardResponse(c, http.StatusUnauthorized, "Invalid token signature", nil, nil)
				}
				return helpers.StandardResponse(c, http.StatusBadRequest, "Bad request", nil, nil)
			}
			if !token.Valid {
				return helpers.StandardResponse(c, http.StatusUnauthorized, "Invalid token", nil, nil)
			}
			if claims.UserID <= 0 {
				return helpers.StandardResponse(c, http.StatusUnauthorized, "user_id missing in token", nil, nil)
			}
			if !roleAllowed(claims.Roles) {
				return helpers.StandardResponse(c, http.StatusForbidden, "Forbidden: required role not present", nil, nil)
			}

			// simpan ke context
			c.Set("claims", claims)
			c.Set(ContextUserIDKey, claims.UserID)
			c.Set("roles", claims.Roles)

			return next(c)
		}
	}
}

// Helper agar controller gampang ambil userID
func CurrentUserID(c echo.Context) (int, error) {
	v := c.Get(ContextUserIDKey)
	if v == nil {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "no authenticated user in context")
	}
	id, ok := v.(int)
	if !ok || id <= 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid user id in context")
	}
	return id, nil
}
