package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"net/http"
)

type aclRequest struct {
	Access   user.AccessType `form:"access"`
	Token    string          `form:"token"`
	Username string          `from:"username"`
	Password string          `form:"password"`
	Topic    string          `form:"topic"`
}

func ACL(ctx *gin.Context) {
	request := &aclRequest{}
	err := ctx.ShouldBind(request)
	if err != nil {
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
		ctx.String(http.StatusUnauthorized, "request is not authorized")
		return
	}
	ctx.String(http.StatusOK, "ok")
}
