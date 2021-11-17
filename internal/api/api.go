package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ReSTServer will return fiber app.
func ReSTServer() *fiber.App {
	app := fiber.New()

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	prometheus := fiberprometheus.NewWith("http", "dispatching", "soteria")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Post("/auth", Auth)
	app.Post("/acl", ACL)

	return app
}
