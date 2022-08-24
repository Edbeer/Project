package utils

import (
	"context"
	"net/http"

	"github.com/Edbeer/Project/config"
	"github.com/labstack/echo/v4"
)

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// ReqIDCtxKey is a key used for the Request ID from context
type ReqIDCtxKey struct{}

// Get user IP address
func GetIP(c echo.Context) string {
	return c.Request().RemoteAddr
}

// Get context  with request id
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

// Configure JWT cookie
func ConfigureJWTCookie(cfg *config.Config, jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:       cfg.Cookie.Name,
		Value:      jwtToken,
		Path:       "/",
		// Domain:     "/",
		RawExpires: "",
		MaxAge:     cfg.Cookie.MaxAge,
		Secure:     cfg.Cookie.Secure,
		// it is better to keep jwtokens on httpOnly from XSS attacks
		HttpOnly:   cfg.Cookie.HTTPOnly, // true
		SameSite:   0,
	}
}

func DeleteCookie(c echo.Context, cookieName string) {
	c.SetCookie(&http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}