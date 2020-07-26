package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitRestApi(port int) {
	router := gin.Default()

	accounts := router.Group("/accounts")
	{
		accounts.POST("/", UserCreate)
		accounts.GET("/:username", UserRead)
		accounts.PUT("/:username", UserUpdate)
		accounts.DELETE("/:username", UserDelete)
	}

	router.POST("/auth", Auth)
	router.POST("/acl", ACL)
	router.POST("/token", Token)

	router.Run(fmt.Sprintf(":%v", port))
}
