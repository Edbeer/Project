package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/pkg/httpe"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")
			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					return c.JSON(httpe.ErrorResponse(httpe.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := validateJWTToken(tokenString, mw.user, c, mw.config); err != nil {
					return c.JSON(httpe.ErrorResponse(err))
				}
				return next(c)
			} else {
				cookie, err := c.Cookie("jwt-token")
				if err != nil {
					return c.JSON(httpe.ErrorResponse(err))
				}

				if err := validateJWTToken(cookie.Value, mw.user, c, mw.config); err != nil {
					return c.JSON(http.StatusUnauthorized, httpe.NewUnauthorizedError(httpe.Unauthorized))
				}
				return next(c)
			}
		}
	}
}

func validateJWTToken(tokenString string, user UserService, c echo.Context, config *config.Config) error {
	if tokenString == "" {
		return httpe.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}
		secret := []byte(config.Server.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return httpe.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpe.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := user.GetUserByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), "user", u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}