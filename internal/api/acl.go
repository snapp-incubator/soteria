package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// aclRequest is the body payload structure of the ACL endpoint.
type aclRequest struct {
	Access   acl.AccessType `form:"access"`
	Token    string         `form:"token"`
	Username string         `from:"username"`
	Password string         `form:"password"`
	Topic    string         `form:"topic"`
}

// ACL is the handler responsible for ACL requests.
// nolint: wrapcheck, funlen
func (a API) ACL(c *fiber.Ctx) error {
	_, span := a.Tracer.Start(c.Context(), "api.acl")
	defer span.End()

	request := new(aclRequest)
	if err := c.BodyParser(request); err != nil {
		a.Logger.
			Warn("acl bad request",
				zap.Error(err),
				zap.String("access", request.Access.String()),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	tokenString := request.Token

	if len(request.Token) == 0 {
		tokenString = request.Username
	}

	if len(tokenString) == 0 {
		tokenString = request.Password
	}

	span.SetAttributes(
		attribute.String("access", request.Access.String()),
		attribute.String("topic", request.Topic),
		attribute.String("token", request.Token),
		attribute.String("username", request.Password),
		attribute.String("password", request.Username),
	)

	topic := request.Topic

	ok, err := a.Authenticator(request.Password).ACL(request.Access, tokenString, topic)
	if err != nil || !ok {
		if err != nil {
			span.RecordError(err)
		}

		var tnaErr authenticator.TopicNotAllowedError

		if errors.As(err, &tnaErr) {
			a.Logger.
				Warn("acl request is not authorized",
					zap.Error(tnaErr))
		} else {
			a.Logger.
				Error("acl request is not authorized",
					zap.Error(err))
		}

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	a.Logger.
		Info("acl ok",
			zap.String("access", request.Access.String()),
			zap.String("topic", request.Topic),
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	return c.Status(http.StatusOK).SendString("ok")
}
