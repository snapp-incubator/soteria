package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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

	tokenString := request.Token
	if len(tokenString) == 0 {
		tokenString = request.Username
	}

	if len(tokenString) == 0 {
		tokenString = request.Password
	}

	span.SetAttributes(
		attribute.String("token", request.Token),
		attribute.String("username", request.Password),
		attribute.String("password", request.Username),
	)

	auth, ok := a.Authenticator(request.Password)
	if !ok {
		a.Logger.Warn("vendor not supported", zap.String("vendor", request.Password))

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	if err := auth.Auth(tokenString); err != nil {
		span.RecordError(err)

		a.Logger.
			Error("auth request is not authorized",
				zap.Error(err),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	a.Logger.
		Info("auth ok",
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	return c.Status(http.StatusOK).SendString("ok")
}
