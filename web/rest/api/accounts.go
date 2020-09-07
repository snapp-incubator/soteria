package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	accountsInfo "gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
)

// Response is the response structure of the REST API
type Response struct {
	Code    accountsInfo.Code `json:"code"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data"`
}

// CreateResponse returns a HTTP Status Code and a response
func CreateResponse(code accountsInfo.Code, data interface{}, details ...string) (int, *Response) {
	return code.HttpStatusCode(), &Response{
		Code:    code,
		Message: fmt.Sprintf("%s: %s", code.Message(), details),
		Data:    data,
	}
}

// createAccountPayload is the body payload structure of create account endpoint
type createAccountPayload struct {
	Username string        `json:"username" form:"username" binding:"required"`
	Password string        `json:"password" form:"password" binding:"required"`
	UserType user.UserType `json:"user_type" form:"user_type" binding:"required"`
}

// CreateAccount is the handler of the create account endpoint
func CreateAccount(ctx *gin.Context) {
	var p createAccountPayload
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(CreateResponse(accountsInfo.BadRequestPayload, nil, err.Error()))
		return
	}

	if err := app.GetInstance().AccountsService.SignUp(p.Username, p.Password, p.UserType); err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}

// ReadAccount is the handler of the read account endpoint
func ReadAccount(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	password := ctx.MustGet("password").(string)

	u, err := app.GetInstance().AccountsService.Info(username, password)
	if err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, u))
}

// updateAccountPayload is the body payload structure of update account endpoint
type updateAccountPayload struct {
	NewPassword string   `json:"new_password" form:"new_password"`
	Secret      string   `json:"secret" form:"secret"`
	IPs         []string `json:"ips" form:"ips"`
}

// UpdateAccount is the handler of the update account endpoint
func UpdateAccount(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	password := ctx.MustGet("password").(string)

	var p updateAccountPayload
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(CreateResponse(accountsInfo.BadRequestPayload, nil, err.Error()))
		return
	}

	err := app.GetInstance().AccountsService.Update(username, password, p.NewPassword, p.Secret, p.IPs)
	if err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}

// DeleteAccount is the handler of the delete account endpoint
func DeleteAccount(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)
	password := ctx.MustGet("password").(string)

	err := app.GetInstance().AccountsService.Delete(username, password)
	if err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}
