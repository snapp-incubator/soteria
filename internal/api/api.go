package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"go.uber.org/zap"
)

type API struct {
	App app.App
}

// ReSTServer will return fiber app.
func (a API) ReSTServer() *fiber.App {
	app := fiber.New()

	// nolint: exhaustivestruct
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
	}))

	prometheus := fiberprometheus.NewWith("http", "dispatching", "soteria")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Post("/auth", a.Auth)
	app.Post("/acl", a.ACL)

	return app
}
