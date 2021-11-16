package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ReSTServer will return an HTTP.Server with given port and gin mode
func ReSTServer() *fiber.App {
	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	app.Post("/auth", Auth)
	app.Post("/acl", ACL)

	prometheus := fiberprometheus.NewWith("soteria", "dispatching", "http")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	return app
}
