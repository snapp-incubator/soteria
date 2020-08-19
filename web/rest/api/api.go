package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func RestServer(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupRouter(),
	}
}
