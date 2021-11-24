package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"go.uber.org/zap"
)

const EMQAuthIgnore = "ignore"

// authRequest is the body payload structure of the auth endpoint.
type authRequest struct {
	Token    string `form:"token"`
	Username string `from:"username"`
	Password string `form:"password"`
}

// Auth is the handler responsible for authentication.
func Auth(c *fiber.Ctx) error {
	_, span := app.GetInstance().Tracer.Start(c.Context(), "api.auth")
	defer span.End()

	s := time.Now()
	request := new(authRequest)

	if err := c.BodyParser(request); err != nil {
		zap.L().
			Warn("bad request",
				zap.Error(err),
			)
		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	tokenString := request.Token
	if len(tokenString) == 0 {
		tokenString = request.Username
	}

	if len(tokenString) == 0 {
		tokenString = request.Password
	}

	if err := app.GetInstance().Authenticator.Auth(tokenString); err != nil {
		span.RecordError(err)

		zap.L().
			Error("auth request is not authorized",
				zap.Error(err),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	zap.L().
		Info("auth ok",
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))

	return c.Status(http.StatusOK).SendString("ok")
}
