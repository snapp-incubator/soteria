package api

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// setupRouter will attach all routes needed for Soteria to gin's default router
func setupRouter(mode string) *gin.Engine {
	gin.SetMode(mode)

	router := gin.New()
	router.Use(gin.Recovery())

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

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}

// RestServer will return an HTTP.Server with given port and gin mode
func RestServer(mode string, port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: setupRouter(mode),
	}
}
