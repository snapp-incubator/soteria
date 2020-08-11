package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	a := router.Group("/accounts")
	{
		a.POST("", CreateAccount)

		authorizedRoutes := a.Use(accountsBasicAuth())
		{
			authorizedRoutes.GET("/:username", ReadAccount)
			authorizedRoutes.PUT("/:username", UpdateAccount)
			authorizedRoutes.DELETE("/:username", DeleteAccount)
		}
	}

	router.POST("/auth", Auth)
	router.POST("/acl", ACL)
	router.POST("/token", Token)

	return router
}

func RunRestApi(port int) error {
	router := setupRouter()
	return router.Run(fmt.Sprintf(":%d", port))
}
