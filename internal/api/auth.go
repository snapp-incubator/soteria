package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// AuthRequest is the body payload structure of the auth endpoint.
type AuthRequest struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

type AuthResponse struct {
	Result      string `json:"result,omitempty"`
	IsSuperuser bool   `json:"is_superuser,omitempty"`
	ExpireAt    int64  `json:"expire_at,omitempty"`
}

// Auth is the handler responsible for authentication.
// Endpoint will be used by EMQ version 5 which supports JSON on both request and response.
// https://www.emqx.io/docs/en/latest/access-control/authn/http.html
// nolint: funlen
func (a API) Authv2(c *fiber.Ctx) error {
	ctx, span := a.Tracer.Start(c.Context(), "api.v2.auth")
	defer span.End()

	request := new(AuthRequest)

	if err := c.BodyParser(request); err != nil {
		span.RecordError(err)

		a.Logger.
			Warn("bad request",
				zap.Error(err),
			)
		authenticator.IncrementWithErrorAuthCounter("unknown_company_before_parse_body", err)

		return c.Status(http.StatusOK).JSON(AuthResponse{
			Result:      "deny",
			IsSuperuser: false,
			ExpireAt:    0,
		})
	}

	vendor, token := ExtractVendorToken(request.Token, request.Username, request.Password)

	auth := a.Authenticator(vendor)

	span.SetAttributes(
		attribute.String("authenticator", auth.GetCompany()),
		attribute.String("cliend-id", request.ClientID),
	)

	if err := auth.Auth(ctx, token); err != nil {
		span.RecordError(err)
		authenticator.IncrementWithErrorAuthCounter(vendor, err)

		if !errors.Is(err, jwt.ErrTokenExpired) {
			a.Logger.
				Error("auth request is not authorized",
					zap.Error(err),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", auth.GetCompany()),
				)
		}

		return c.Status(http.StatusOK).JSON(AuthResponse{
			Result:      "deny",
			IsSuperuser: false,
			ExpireAt:    0,
		})
	}

	a.Logger.
		Info("auth ok",
			zap.String("token", request.Token),
			zap.String("username", request.Username),
			zap.String("password", request.Password),
			zap.String("authenticator", auth.GetCompany()),
		)
	authenticator.IncrementAuthCounter(vendor)

	return c.Status(http.StatusOK).JSON(AuthResponse{
		Result:      "allow",
		IsSuperuser: auth.IsSuperuser(),
		ExpireAt:    0,
	})
}
