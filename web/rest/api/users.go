package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserCreate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

func UserRead(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

func UserUpdate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}

func UserDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status": "200",
	})
}
