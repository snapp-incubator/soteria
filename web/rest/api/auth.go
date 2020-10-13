package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// authRequest is the body payload structure of the auth endpoint
type authRequest struct {
	Token    string `form:"token"`
	Username string `from:"username"`
	Password string `form:"password"`
}

// Auth is the handler responsible for authentication
func Auth(ctx *gin.Context) {
	s := time.Now()
	request := &authRequest{}
	err := ctx.ShouldBind(request)
	if err != nil {
		zap.L().
			Warn("bad request",
				zap.Error(err),
			)
		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))
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
	ok, err := app.GetInstance().Authenticator.Auth(tokenString)
	if err != nil || !ok {

		zap.L().
			Error("auth request is not authorized",
				zap.Error(err),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")
		return
	}

	zap.L().
		Info("auth ok",
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.HttpApi, internal.Soteria, internal.Auth, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.HttpApi, internal.Soteria, internal.Auth, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.HttpApi, internal.Soteria, internal.Auth, float64(time.Since(s).Nanoseconds()))
	ctx.String(http.StatusOK, "ok")
}
