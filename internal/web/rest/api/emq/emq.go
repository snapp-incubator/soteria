package emq

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/emq"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/web/rest/api/emq/request"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/web/rest/api/emq/response"
)

func Register(group *gin.RouterGroup) {
	group.POST("", Create)
	group.GET("", List)
}

// List returns the list of registered users for emqx-redis-auth.
func List(ctx *gin.Context) {
	users, err := app.GetInstance().EMQStore.LoadAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	ctx.JSON(http.StatusCreated, response.Response{
		Message: "success",
		Data:    users,
	})
}

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

	token, err := app.GetInstance().Authenticator.SuperuserToken(p.Username, time.Duration(p.Duration))
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
		Data: response.Create{
			Token:    token,
			Password: p.Password,
		},
	})
}
