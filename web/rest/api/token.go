package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Token(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}
