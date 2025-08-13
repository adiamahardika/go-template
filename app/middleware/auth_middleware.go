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

func JWTAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return helpers.StandardResponse(c, http.StatusUnauthorized, "Missing Authorization header", nil, nil)
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
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

			// Check for admin role
			isAdmin := false
			for _, role := range claims.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				return helpers.StandardResponse(c, http.StatusForbidden, "Forbidden: Admin access required", nil, nil)
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}