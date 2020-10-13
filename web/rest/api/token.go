package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
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
	s := time.Now()
	request := &TokenRequest{}
	err := ctx.Bind(request)
	if err != nil {
		zap.L().
			Warn("bad request",
				zap.Error(err),
			)
		app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Token, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}
	tokenString, err := app.GetInstance().Authenticator.Token(request.GrantType, request.ClientID, request.ClientSecret)
	if err != nil {

		zap.L().
			Error("token request is not authorized",
				zap.Error(err),
				zap.String("grant type", string(request.GrantType)),
				zap.String("client id", request.ClientID),
				zap.String("client secret", request.ClientSecret),
			)

		if errors.Is(err, db.ErrDb) {
			app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Token, http.StatusInternalServerError)
			app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, "database error happened")
			app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, request.ClientID)
			app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
			ctx.String(http.StatusInternalServerError, "internal server error")
			return
		}

		app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Token, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, request.ClientID)
		app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")
		return
	}

	zap.L().
		Info("token request accepted",
			zap.String("grant type", string(request.GrantType)),
			zap.String("client id", request.ClientID),
			zap.String("client secret", request.ClientSecret),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Token, http.StatusAccepted)
	app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Success, "token request accepted")
	app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Success, request.ClientID)
	app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
	ctx.String(http.StatusAccepted, tokenString)
}
