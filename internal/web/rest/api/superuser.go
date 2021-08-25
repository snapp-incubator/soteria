package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"go.uber.org/zap"
)

// Superuser is the handler responsible for authentication a superuser.
// it uses the exact same request payload as authentication endpoint.
// emq calls the superuser endpoint after success auth based on:
// https://github.com/emqx/emqx-auth-http/blob/master/src/emqx_auth_http.erl#L45.
func Superuser(ctx *gin.Context) {
	authSpan := app.GetInstance().Tracer.StartSpan("api.rest.auth")
	defer authSpan.Finish()

	s := time.Now()

	var request authRequest

	if err := ctx.ShouldBind(&request); err != nil {
		zap.L().
			Warn("bad request",
				zap.Error(err),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Superuser, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Superuser, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Superuser, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}

	tokenString := request.Token

	if len(tokenString) == 0 {
		tokenString = request.Username
	}

	if len(tokenString) == 0 {
		tokenString = request.Password
	}

	authCheckSpan := app.GetInstance().Tracer.StartSpan("auth check", opentracing.ChildOf(authSpan.Context()))
	defer authCheckSpan.Finish()

	iss, err := app.GetInstance().Authenticator.Issuer(ctx, tokenString)
	if err != nil {
		authCheckSpan.SetTag("success", false)
		authCheckSpan.SetTag("error", err.Error())

		zap.L().
			Error("superuser request is not authorized",
				zap.Error(err),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Superuser, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Superuser, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Superuser, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")

		return
	}

	if iss != user.ThirdParty {
		zap.L().
			Error("superuser request is not authorized",
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Superuser, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Superuser, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Superuser, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")

		return
	}

	if _, err := app.GetInstance().Authenticator.Auth(ctx, tokenString); err != nil {
		authCheckSpan.SetTag("success", false)
		authCheckSpan.SetTag("error", err.Error())

		zap.L().
			Error("superuser request is not authorized",
				zap.Error(err),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Superuser, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Superuser, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Superuser, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")

		return
	}

	zap.L().
		Info("superuser ok",
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Superuser, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Superuser, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Superuser, float64(time.Since(s).Nanoseconds()))

	authSpan.SetTag("success", true)

	ctx.String(http.StatusOK, "ok")
}
