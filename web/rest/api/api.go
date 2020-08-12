package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
)

type Core struct {
	Authenticator *accounts.Authenticator
}

func setupRouter(c *Core) *gin.Engine {

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

	router.POST("/auth", c.Auth)
	router.POST("/acl", c.ACL)
	router.POST("/token", c.Token)

	return router
}

func RunRestApi(c *Core, port int) error {
	router := setupRouter(c)
	return router.Run(fmt.Sprintf(":%d", port))
}
