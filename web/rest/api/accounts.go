package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	accountsInfo "gitlab.snapp.ir/dispatching/soteria/pkg/errors"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"time"
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
	NewPassword     string        `json:"new_password" form:"new_password"`
	IPs             []string      `json:"ips" form:"ips"`
	Secret          string        `json:"secret" form:"secret"`
	Type            user.UserType `json:"type" form:"type"`
	TokenExpiration time.Duration `json:"token_expiration" form:"token_expiration"`
}

// UpdateAccount is the handler of the update account endpoint
func UpdateAccount(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	var p updateAccountPayload
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(CreateResponse(accountsInfo.BadRequestPayload, nil, err.Error()))
		return
	}

	if err := app.GetInstance().AccountsService.Update(username, p.NewPassword, p.Type, p.IPs, p.Secret, p.TokenExpiration); err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}

// DeleteAccount is the handler of the delete account endpoint
func DeleteAccount(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	if err := app.GetInstance().AccountsService.Delete(username); err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}

// createAccountRulePayload is the body payload structure of create account rule endpoint
type createAccountRulePayload struct {
	Endpoint   string         `json:"endpoint" form:"endpoint"`
	Topic      topics.Type    `json:"topic" form:"topic"`
	AccessType acl.AccessType `json:"access_type" form:"access_type"`
}

// CreateAccountRule is the handler of the create account rule endpoint
func CreateAccountRule(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	var p createAccountRulePayload
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(CreateResponse(accountsInfo.BadRequestPayload, nil, err.Error()))
		return
	}

	r, err := app.GetInstance().AccountsService.CreateRule(username, p.Endpoint, p.Topic, p.AccessType)
	if err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, r))
}

// ReadAccountRule is the handler of the read account rule endpoint
func ReadAccountRule(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	plainUUID := ctx.Param("uuid")
	if plainUUID == "" {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil))
		return
	}

	ruleUUID, err := uuid.Parse(plainUUID)
	if err != nil {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil, err.Error()))
		return
	}

	r, rErr := app.GetInstance().AccountsService.GetRule(username, ruleUUID)
	if rErr != nil {
		ctx.JSON(CreateResponse(rErr.Code, nil, rErr.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, r))
}

// updateAccountRulePayload is the body payload structure of update account rule endpoint
type updateAccountRulePayload struct {
	Endpoint   string         `json:"endpoint" form:"endpoint"`
	Topic      topics.Type    `json:"topic" form:"topic"`
	AccessType acl.AccessType `json:"access_type" form:"access_type"`
}

// UpdateAccountRule is the handler of the update account rule endpoint
func UpdateAccountRule(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	var p updateAccountRulePayload
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(CreateResponse(accountsInfo.BadRequestPayload, nil, err.Error()))
		return
	}

	plainUUID := ctx.Param("uuid")
	if plainUUID == "" {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil))
		return
	}

	ruleUUID, err := uuid.Parse(plainUUID)
	if err != nil {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil, err.Error()))
		return
	}

	if err := app.GetInstance().AccountsService.UpdateRule(username, ruleUUID, p.Endpoint, p.Topic, p.AccessType); err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}

// DeleteAccountRule is the handler of the delete account rule endpoint
func DeleteAccountRule(ctx *gin.Context) {
	username := ctx.MustGet("username").(string)

	plainUUID := ctx.Param("uuid")
	if plainUUID == "" {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil))
		return
	}

	ruleUUID, err := uuid.Parse(plainUUID)
	if err != nil {
		ctx.JSON(CreateResponse(accountsInfo.InvalidRuleUUID, nil, err.Error()))
		return
	}

	if err := app.GetInstance().AccountsService.DeleteRule(username, ruleUUID); err != nil {
		ctx.JSON(CreateResponse(err.Code, nil, err.Message))
		return
	}

	ctx.JSON(CreateResponse(accountsInfo.SuccessfulOperation, nil))
}
