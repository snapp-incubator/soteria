package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}
