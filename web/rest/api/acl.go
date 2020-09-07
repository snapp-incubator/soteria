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

// aclRequest is the body payload structure of the ACL endpoint
type aclRequest struct {
	Access   user.AccessType `form:"access"`
	Token    string          `form:"token"`
	Username string          `from:"username"`
	Password string          `form:"password"`
	Topic    string          `form:"topic"`
}

// ACL is the handler responsible for ACL requests
func ACL(ctx *gin.Context) {
	s := time.Now()
	request := &aclRequest{}
	err := ctx.ShouldBind(request)
	if err != nil {

		zap.L().
			Warn("acl bad request",
				zap.Error(err),
				zap.String("access", string(request.Access)),
				zap.String("topic", request.Topic),
				zap.String("token", request.Token),
				zap.String("username", request.Password),
				zap.String("password", request.Username),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Acl, http.StatusBadRequest)
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Acl, internal.Failure, "bad request")
		app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
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
	ok, err := app.GetInstance().Authenticator.Acl(request.Access, tokenString, request.Topic)
	if err != nil || !ok {

		zap.L().
			Error("acl request is not authorized",
				zap.Error(err),
			)

		app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Acl, http.StatusUnauthorized)
		app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Acl, internal.Failure, "request is not authorized")
		app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
		ctx.String(http.StatusUnauthorized, "request is not authorized")
		return
	}

	zap.L().
		Info("acl ok",
			zap.String("access", string(request.Access)),
			zap.String("topic", request.Topic),
			zap.String("token", request.Token),
			zap.String("username", request.Password),
			zap.String("password", request.Username),
		)

	app.GetInstance().Metrics.ObserveStatusCode(internal.Soteria, internal.Acl, http.StatusOK)
	app.GetInstance().Metrics.ObserveStatus(internal.Soteria, internal.Acl, internal.Success, "ok")
	app.GetInstance().Metrics.ObserveResponseTime(internal.Soteria, internal.Acl, float64(time.Since(s).Nanoseconds()))
	ctx.String(http.StatusOK, "ok")
}
