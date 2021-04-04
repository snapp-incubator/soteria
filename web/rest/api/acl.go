package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// aclRequest is the body payload structure of the ACL endpoint
type aclRequest struct {
	Access   acl.AccessType `form:"access"`
	Token    string         `form:"token"`
	Username string         `from:"username"`
	Password string         `form:"password"`
	Topic    string         `form:"topic"`
}

// ACL is the handler responsible for ACL requests
func ACL(ctx *gin.Context) {
	aclSpan := app.GetInstance().Tracer.StartSpan("api.rest.acl")
	defer aclSpan.Finish()

	s := time.Now()
	request := &aclRequest{}
	err := ctx.ShouldBind(request)
	if err != nil {

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
		ctx.String(http.StatusBadRequest, "bad request")
		return
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
		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}

	aclCheckSpan := app.GetInstance().Tracer.StartSpan("acl check", opentracing.ChildOf(aclSpan.Context()))

	ok, err := app.GetInstance().Authenticator.Acl(ctx, request.Access, tokenString, topic)
	if err != nil || !ok {
		aclCheckSpan.SetTag("success", false)
		if err != nil {
			aclCheckSpan.SetTag("error", err.Error())
		}

		if errors.Is(err, authenticator.TopicNotAllowed) {
			zap.L().
				Warn("acl request is not authorized",
					zap.Error(err))
		} else {
			zap.L().
				Error("acl request is not authorized",
					zap.Error(err))
		}

		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, request.Access.String(), internal.Failure, string(topicType))
			app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
			ctx.String(http.StatusInternalServerError, "internal server error")

			aclCheckSpan.Finish()

			return
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Acl, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, request.Access.String(), internal.Failure, string(topicType))
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Acl, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")

		aclCheckSpan.Finish()

		return
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

	ctx.String(http.StatusOK, "ok")
}
