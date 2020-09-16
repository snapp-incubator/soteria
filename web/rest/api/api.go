package api

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// setupRouter will attach all routes needed for Soteria to gin's default router
func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, false))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))

	a := router.Group("/accounts")
	{
		a.POST("", CreateAccount)
		authorizedRoutes := a.Use(accountsBasicAuth())
		{
			authorizedRoutes.GET("/:username", ReadAccount)
			authorizedRoutes.PUT("/:username", UpdateAccount)
			authorizedRoutes.DELETE("/:username", DeleteAccount)

			authorizedRoutes.POST("/:username/rules", CreateAccountRule)
			authorizedRoutes.GET("/:username/rules/:uuid", ReadAccountRule)
			authorizedRoutes.PUT("/:username/rules/:uuid", UpdateAccountRule)
			authorizedRoutes.DELETE("/:username/rules/:uuid", DeleteAccountRule)
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
