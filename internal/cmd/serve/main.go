package serve

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
	"github.com/snapp-incubator/soteria/internal/api"
	"github.com/snapp-incubator/soteria/internal/authenticator"
	"github.com/snapp-incubator/soteria/internal/clientid"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/metric"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Serve struct {
	Cfg    config.Config
	Logger *zap.Logger
	Tracer trace.Tracer
}

// Register serve command.
func (s Serve) Register(root *cobra.Command) {
	root.AddCommand(
		//nolint: exhaustruct
		&cobra.Command{
			Use:   "serve",
			Short: "serve runs the application",
			Long:  `serve will run Soteria ReST server and waits until user disrupts.`,
			Run: func(_ *cobra.Command, _ []string) {
				s.main()
			},
		},
	)
}

func (s Serve) main() {
	auth, err := authenticator.Builder{
		Vendors:         s.Cfg.Vendors,
		Logger:          s.Logger,
		ValidatorConfig: s.Cfg.Validator,
		Tracer:          s.Tracer,
	}.Authenticators()
	if err != nil {
		s.Logger.Fatal("authenticator building failed", zap.Error(err))
	}

	api := api.API{
		DefaultVendor:  s.Cfg.DefaultVendor,
		Authenticators: auth,
		Tracer:         s.Tracer,
		Logger:         s.Logger.Named("api"),
		Parser:         clientid.NewParser(s.Cfg.Parser),
		Metrics:        metric.NewAPIMetrics(),
	}

	if _, ok := api.Authenticators[s.Cfg.DefaultVendor]; !ok {
		s.Logger.Fatal("default vendor shouldn't be nil, please set it")
	}

	rest := api.ReSTServer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//nolint: exhaustruct
	if err := rest.Listen(fmt.Sprintf(":%d", s.Cfg.HTTPPort), fiber.ListenConfig{
		GracefulContext: ctx,
	}); err != nil {
		s.Logger.Fatal("failed to run REST HTTP server", zap.Error(err))
	}
}
