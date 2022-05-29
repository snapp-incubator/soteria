package api

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type API struct {
	Authenticator *authenticator.Authenticator
	Tracer        trace.Tracer
	Logger        zap.Logger
}

// ReSTServer will return fiber app.
func (a API) ReSTServer() *fiber.App {
	app := fiber.New()

	// nolint: exhaustruct
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: a.Logger.Named("fiber"),
	}))

	prometheus := fiberprometheus.NewWith("http", "dispatching", "soteria")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Post("/auth", a.Auth)
	app.Post("/acl", a.ACL)

	return app
}
