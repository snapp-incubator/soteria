package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRequest struct {
	Token    string `form:"token"`
	Username string `from:"username"`
	Password string `form:"password"`
}

func (c *Core) Auth(ctx *gin.Context) {
	request := &authRequest{}
	err := ctx.ShouldBind(request)
	if err != nil {
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
	if len(tokenString) == 0 {
		ctx.String(http.StatusBadRequest, "")
	}
	ok, err := c.Authenticator.Auth(tokenString)
	if err != nil || !ok {
		ctx.String(http.StatusUnauthorized, "request is not authorized")
		return
	}
	ctx.String(http.StatusOK, "ok")
}
