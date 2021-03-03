package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/acl"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TokenRequest is the body payload structure of the token endpoint
type TokenRequest struct {
	GrantType    acl.AccessType `json:"grant_type" form:"grant_type" query:"grant_type"`
	ClientID     string         `json:"client_id" form:"client_id" query:"client_id"`
	ClientSecret string         `json:"client_secret" form:"client_secret" query:"client_secret"`
}

// Token is the handler responsible for Token requests
func Token(ctx *gin.Context) {
	tokenSpan := app.GetInstance().Tracer.StartSpan("api.rest.token")
	defer tokenSpan.Finish()

	s := time.Now()
	request := &TokenRequest{}
	err := ctx.Bind(request)
	if err != nil {
		zap.L().
			Warn("bad request",
				zap.Error(err),
			)
		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Token, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}

	tokenIssueSpan := app.GetInstance().Tracer.StartSpan("issue token", opentracing.ChildOf(tokenSpan.Context()))

	tokenString, err := app.GetInstance().Authenticator.Token(ctx, request.GrantType, request.ClientID, request.ClientSecret)
	if err != nil {
		zap.L().
			Error("token request is not authorized",
				zap.Error(err),
				zap.String("grant_type", request.GrantType.String()),
				zap.String("client_id", request.ClientID),
				zap.String("client_secret", request.ClientSecret),
			)

		tokenIssueSpan.SetTag("success", false)

		tokenIssueSpan.SetTag("error", err.Error()).Finish()

		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Token, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Failure, request.ClientID)
			app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
			ctx.String(http.StatusInternalServerError, "internal server error")

			tokenIssueSpan.Finish()

			return
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Token, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Failure, request.ClientID)
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")

		tokenIssueSpan.Finish()

		return
	}

	zap.L().
		Info("token request accepted",
			zap.String("grant_type", request.GrantType.String()),
			zap.String("client_id", request.ClientID),
			zap.String("client_secret", request.ClientSecret),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Token, http.StatusAccepted)
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Success, "token request accepted")
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Token, internal.Success, request.ClientID)
	app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))

	tokenIssueSpan.SetTag("success", true).Finish()

	ctx.String(http.StatusAccepted, tokenString)
}
