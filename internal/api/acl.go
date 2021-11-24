package api

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
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
func ACL(c *fiber.Ctx) error {
	_, span := app.GetInstance().Tracer.Start(c.Context(), "api.acl")
	defer span.End()

	request := new(aclRequest)
	if err := c.BodyParser(request); err != nil {
		zap.L().
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

	topic := topics.Topic(request.Topic)
	topicType := topic.GetType()

	if len(topicType) == 0 {
		zap.L().
			Warn("acl bad request",
				zap.String("access", request.Access.String()),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	ok, err := app.GetInstance().Authenticator.ACL(request.Access, tokenString, topic)
	if err != nil || !ok {
		if err != nil {
			span.RecordError(err)
		}

		if errors.Is(err, authenticator.TopicNotAllowedError{}) {
			zap.L().
				Warn("acl request is not authorized",
					zap.Error(err))
		} else {
			zap.L().
				Error("acl request is not authorized",
					zap.Error(err))
		}

		return c.Status(http.StatusUnauthorized).SendString("request is not authorized")
	}

	zap.L().
		Info("acl ok",
			zap.String("access", request.Access.String()),
			zap.String("topic", request.Topic),
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	return c.Status(http.StatusOK).SendString("ok")
}
