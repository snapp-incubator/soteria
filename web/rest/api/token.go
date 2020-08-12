package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TokenRequest is payload structure for token request
type TokenRequest struct {
	GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
}

func (c *Core) Token(ctx *gin.Context) {
	request := &TokenRequest{}
	err := ctx.Bind(request)
	if err != nil {
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}
	tokenString, err := c.Authenticator.Token(request.GrantType, request.ClientID, request.ClientSecret)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "unauthorized request")
		return
	}
	ctx.String(http.StatusAccepted, tokenString)
}
