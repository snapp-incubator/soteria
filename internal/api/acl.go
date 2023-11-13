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

type AclResponse struct {
	Result string `json:"result,omitempty"`
}

// AclRequest is the body payload structure of the ACL endpoint.
type AclRequest struct {
	Access   acl.AccessType `form:"access"`
	Token    string         `form:"token"`
	Username string         `from:"username"`
	Password string         `form:"password"`
	Topic    string         `form:"topic"`
}

// ACL is the handler responsible for ACL requests.
// nolint: wrapcheck, funlen
func (a API) ACLv1(c *fiber.Ctx) error {
	_, span := a.Tracer.Start(c.Context(), "api.v1.acl")
	defer span.End()

	request := new(AclRequest)
	if err := c.BodyParser(request); err != nil {
		a.Logger.
			Warn("acl bad request",
				zap.Error(err),
				zap.String("access", request.Access.String()),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Username),
				zap.String("password", request.Password),
			)

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	vendor, token := ExtractVendorToken(request.Token, request.Username, request.Password)

	topic := request.Topic
	auth := a.Authenticator(vendor)

	span.SetAttributes(
		attribute.String("access", request.Access.String()),
		attribute.String("topic", request.Topic),
		attribute.String("token", request.Token),
		attribute.String("username", request.Username),
		attribute.String("password", request.Password),
		attribute.String("authenticator", auth.GetCompany()),
	)

	ok, err := auth.ACL(request.Access, token, topic)
	if err != nil || !ok {
		if err != nil {
			span.RecordError(err)
		}

		var tnaErr authenticator.TopicNotAllowedError

		if errors.As(err, &tnaErr) {
			a.Logger.
				Warn("acl request is not authorized",
					zap.Error(tnaErr),
					zap.String("access", request.Access.String()),
					zap.String("topic", request.Topic),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", auth.GetCompany()))
		} else if !errors.Is(err, jwt.ErrTokenExpired) {
			a.Logger.
				Error("acl request is not authorized",
					zap.Error(err),
					zap.String("access", request.Access.String()),
					zap.String("topic", request.Topic),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", auth.GetCompany()))
		}

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	a.Logger.
		Info("acl ok",
			zap.String("access", request.Access.String()),
			zap.String("topic", request.Topic),
			zap.String("token", request.Token),
			zap.String("username", request.Username),
			zap.String("password", request.Password),
			zap.String("authenticator", auth.GetCompany()),
		)

	return c.Status(http.StatusOK).SendString("ok")
}

// aclRequest is the body payload structure of the ACL endpoint.
type aclv2Request struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
	Action   string `json:"action"`
}

// ACL is the handler responsible for ACL requests.
// nolint: wrapcheck, funlen
func (a API) ACLv2(c *fiber.Ctx) error {
	_, span := a.Tracer.Start(c.Context(), "api.v2.acl")
	defer span.End()

	request := new(aclv2Request)
	if err := c.BodyParser(request); err != nil {
		a.Logger.
			Warn("acl bad request",
				zap.Error(err),
				zap.String("access", request.Action),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Username),
				zap.String("password", request.Password),
			)

		return c.Status(http.StatusBadRequest).JSON(AclResponse{
			Result: "deny",
		})
	}

	vendor, token := ExtractVendorToken(request.Token, request.Username, request.Password)

	topic := request.Topic
	auth := a.Authenticator(vendor)

	span.SetAttributes(
		attribute.String("access", request.Action),
		attribute.String("topic", request.Topic),
		attribute.String("token", request.Token),
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

	ok, err := auth.ACL(access, token, topic)
	if err != nil || !ok {
		if err != nil {
			span.RecordError(err)
		}

		var tnaErr authenticator.TopicNotAllowedError

		if errors.As(err, &tnaErr) {
			a.Logger.
				Warn("acl request is not authorized",
					zap.Error(tnaErr),
					zap.String("access", request.Action),
					zap.String("topic", request.Topic),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", auth.GetCompany()))
		} else if !errors.Is(err, jwt.ErrTokenExpired) {
			a.Logger.
				Error("acl request is not authorized",
					zap.Error(err),
					zap.String("access", request.Action),
					zap.String("topic", request.Topic),
					zap.String("token", request.Token),
					zap.String("username", request.Username),
					zap.String("password", request.Password),
					zap.String("authenticator", auth.GetCompany()))
		}

		return c.Status(http.StatusUnauthorized).JSON(AclResponse{
			Result: "deny",
		})
	}

	a.Logger.
		Info("acl ok",
			zap.String("access", request.Action),
			zap.String("topic", request.Topic),
			zap.String("token", request.Token),
			zap.String("username", request.Username),
			zap.String("password", request.Password),
			zap.String("authenticator", auth.GetCompany()),
		)

	return c.Status(http.StatusOK).JSON(AclResponse{
		Result: "allow",
	})
}
