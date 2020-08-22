package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TokenRequest is payload structure for token request
type TokenRequest struct {
	GrantType    user.AccessType `json:"grant_type" form:"grant_type" query:"grant_type"`
	ClientID     string          `json:"client_id" form:"client_id" query:"client_id"`
	ClientSecret string          `json:"client_secret" form:"client_secret" query:"client_secret"`
}

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

		app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Token, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Token, internal.Failure, "request is not authorized")
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
	app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Token, float64(time.Since(s).Nanoseconds()))
	ctx.String(http.StatusAccepted, tokenString)
}
