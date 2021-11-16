package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
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
	aclSpan := app.GetInstance().Tracer.StartSpan("api.rest.acl")
	defer aclSpan.Finish()

	s := time.Now()

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

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	tokenString := request.Token

	if len(request.Token) == 0 {
		tokenString = request.Username
	}

	if len(tokenString) == 0 {
		tokenString = request.Password
	}

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

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))

		return c.Status(http.StatusBadRequest).SendString("bad request")
	}

	aclCheckSpan := app.GetInstance().Tracer.StartSpan("acl check", opentracing.ChildOf(aclSpan.Context()))
	defer aclCheckSpan.Finish()

	ok, err := app.GetInstance().Authenticator.ACL(request.Access, tokenString, topic)
	if err != nil || !ok {
		aclCheckSpan.SetTag("success", false)

		if err != nil {
			aclCheckSpan.SetTag("error", err.Error())
		}

		// nolint: exhaustivestruct
		if errors.Is(err, authenticator.TopicNotAllowedError{}) {
			zap.L().
				Warn("acl request is not authorized",
					zap.Error(err))
		} else {
			zap.L().
				Error("acl request is not authorized",
					zap.Error(err))
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, request.Access.String(), internal.Failure, string(topicType))
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))

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

	app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, request.Access.String(), internal.Success, string(topicType))
	app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))

	aclCheckSpan.SetTag("success", true)

	return c.Status(http.StatusOK).SendString("ok")
}
