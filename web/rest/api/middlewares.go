package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	accountsInfo "gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"strings"
)

func accountsBasicAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			ctx.JSON(CreateResponse(accountsInfo.WrongUsernameOrPassword, nil))
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			ctx.JSON(CreateResponse(accountsInfo.WrongUsernameOrPassword, nil))
			return
		}

		if pair[0] != ctx.Param("username") {
			ctx.JSON(CreateResponse(accountsInfo.UsernameMismatch, nil))
			return
		}

		_, err := app.GetInstance().AccountsService.Info(pair[0], pair[1])
		if err != nil {
			ctx.JSON(CreateResponse(err.Code, nil, err.Message))
			return
		}

		ctx.Set("username", pair[0])
		ctx.Set("password", pair[1])

		ctx.Next()
	}
}
