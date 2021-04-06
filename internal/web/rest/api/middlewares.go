package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	accountsInfo "gitlab.snapp.ir/dispatching/soteria/v3/pkg/errors"
	"strings"
)

// accountsBasicAuth is the authentication middleware for the accounts API
func accountsBasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			ctx.AbortWithStatusJSON(CreateResponse(accountsInfo.WrongUsernameOrPassword, nil))
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			ctx.AbortWithStatusJSON(CreateResponse(accountsInfo.WrongUsernameOrPassword, nil))
			return
		}

		if pair[0] != ctx.Param("username") {
			ctx.AbortWithStatusJSON(CreateResponse(accountsInfo.UsernameMismatch, nil))
			return
		}

		_, err := app.GetInstance().AccountsService.Info(ctx, pair[0], pair[1])
		if err != nil {
			ctx.AbortWithStatusJSON(CreateResponse(err.Code, nil, err.Message))
			return
		}

		ctx.Set("username", pair[0])
		ctx.Set("password", pair[1])

		ctx.Next()
	}
}
