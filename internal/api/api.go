package api

import (
	"strings"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const VendorTokenSeparator = ":"

type API struct {
	Authenticators map[string]authenticator.Authenticator
	DefaultVendor  string
	Tracer         trace.Tracer
	Logger         *zap.Logger
}

// MetricLogSkipper check if route is equal "metric" disable log.
func MetricLogSkipper(ctx *fiber.Ctx) bool {
	route := string(ctx.Request().URI().Path())

	return route == "/metrics"
}

// ReSTServer will return fiber app.
func (a API) ReSTServer() *fiber.App {
	app := fiber.New()

	//nolint: exhaustruct
	app.Use(fiberzap.New(fiberzap.Config{
		Next:   MetricLogSkipper,
		Logger: a.Logger.Named("fiber"),
	}))

	prometheus := fiberprometheus.NewWith("http", "dispatching", "soteria")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Post("/auth", a.Authv1)
	app.Post("/acl", a.ACLv1)

	app.Post("/v2/auth", a.Authv2)
	app.Post("/v2/acl", a.ACLv2)

	return app
}

func (a API) Authenticator(vendor string) authenticator.Authenticator {
	auth, ok := a.Authenticators[vendor]

	if ok {
		return auth
	}

	return a.Authenticators[a.DefaultVendor]
}

func ExtractVendorToken(rawToken, username, password string) (string, string) {
	tokenString := rawToken

	if len(tokenString) == 0 {
		tokenString = username
	}

	if len(tokenString) == 0 {
		tokenString = password
	}

	split := strings.Split(tokenString, VendorTokenSeparator)

	var vendor, token string

	if len(split) == 2 { //nolint:mnd
		vendor = split[0]
		token = split[1]
	} else {
		vendor = ""
		token = split[0]
	}

	return vendor, token
}
