package emq

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/emq"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/web/rest/api/emq/request"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/web/rest/api/emq/response"
)

// Create is the handler of the create emq redis account endpoint.
func Create(ctx *gin.Context) {
	var p request.Create

	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	if err := p.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	token, err := app.GetInstance().Authenticator.SuperuserToken(p.Username, p.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	if err := app.GetInstance().EMQStore.Save(ctx, emq.User{
		Username:    token,
		Password:    p.Password,
		IsSuperuser: true,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Message: "success",
		Data:    nil,
	})
}
