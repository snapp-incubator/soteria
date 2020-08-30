package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// setupRouter will attach all routes needed for Soteria to gin's default router
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

// RestServer will return an HTTP.Server with given port and our router
func RestServer(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupRouter(),
	}
}
