package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/pkg/acl"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type ACLResponse struct {
	Result string `json:"result,omitempty"`
}

// ACLRequest is the body payload structure of the ACL endpoint.
type ACLRequest struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
	Action   string `json:"action"`
}

// ACLv2 is the handler responsible for ACL requests coming from EMQv5.
// https://www.emqx.io/docs/en/latest/access-control/authz/http.html
// nolint: funlen
func (a API) ACLv2(c *fiber.Ctx) error {
	ctx, span := a.Tracer.Start(c.Context(), "api.v2.acl")
	defer span.End()

	request := new(ACLRequest)

	err := c.BodyParser(request)
	if err != nil {
		a.Logger.
			Warn("acl bad request",
				zap.Error(err),
				zap.String("access", request.Action),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Username),
				zap.String("password", request.Password),
			)
		a.Metrics.ACLFailed("unknown_company_before_parse_body", err)

		return c.Status(http.StatusOK).JSON(ACLResponse{
			Result: "deny",
		})
	}

	vendor, token := ExtractVendorToken(request.Token, request.Username, request.Password)

	topic := request.Topic
	auth := a.Authenticator(vendor)

	logger := a.Logger.With(
		zap.String("access", request.Action),
		zap.String("topic", request.Topic),
		zap.String("token", request.Token),
		zap.String("username", request.Username),
		zap.String("password", request.Password),
		zap.String("authenticator", auth.GetCompany()),
	)

	span.SetAttributes(
		attribute.String("access", request.Action),
		attribute.String("topic", request.Topic),
		attribute.String("username", request.Username),
		attribute.String("password", request.Password),
		attribute.String("authenticator", auth.GetCompany()),
	)

	var access acl.AccessType

	switch request.Action {
	case "publish":
		access = acl.Pub
	case "subscribe":
		access = acl.Sub
	}

	ok, err := auth.ACL(ctx, access, token, topic)
	if err != nil || !ok {
		if err != nil {
			span.RecordError(err)
		}

		a.Metrics.ACLFailed(auth.GetCompany(), err)

		var tnaErr authenticator.TopicNotAllowedError
		if errors.As(err, &tnaErr) {
			logger.
				Warn("acl request is not authorized",
					zap.Error(tnaErr))
		} else if !errors.Is(err, jwt.ErrTokenExpired) {
			logger.
				Error("acl request is not authorized",
					zap.Error(err))
		}

		return c.Status(http.StatusOK).JSON(ACLResponse{
			Result: "deny",
		})
	}

	logger.
		Info("acl ok")
	a.Metrics.ACLSuccess(auth.GetCompany())

	return c.Status(http.StatusOK).JSON(ACLResponse{
		Result: "allow",
	})
}
