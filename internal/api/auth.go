package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// authRequest is the body payload structure of the auth endpoint.
type authRequest struct {
	Token    string `form:"token"`
	Username string `from:"username"`
	Password string `form:"password"`
}

// Auth is the handler responsible for authentication.
// nolint: wrapcheck
func (a API) Auth(c *fiber.Ctx) error {
	_, span := a.Tracer.Start(c.Context(), "api.auth")
	defer span.End()

	request := new(authRequest)

	if err := c.BodyParser(request); err != nil {
		span.RecordError(err)

		a.Logger.
			Warn("bad request",
				zap.Error(err),
			)

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	vendor, token := ExtractVendorToken(request.Token, request.Username, request.Password)

	authenticator := a.Authenticator(vendor)

	span.SetAttributes(
		attribute.String("token", request.Token),
		attribute.String("username", request.Username),
		attribute.String("password", request.Password),
		attribute.String("authenticator", authenticator.GetCompany()),
	)

	if err := authenticator.Auth(token); err != nil {
		span.RecordError(err)

		if !errors.Is(err, jwt.ErrTokenExpired) {
			a.Logger.
				Error("auth request is not authorized",
					zap.Error(err),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", authenticator.GetCompany()),
				)
		}

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	a.Logger.
		Info("auth ok",
			zap.String("token", request.Token),
			zap.String("username", request.Username),
			zap.String("password", request.Password),
			zap.String("authenticator", authenticator.GetCompany()),
		)

	return c.Status(http.StatusOK).SendString("ok")
}
