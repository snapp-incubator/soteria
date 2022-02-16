package serve

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/api"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Serve struct {
	Cfg    config.Config
	Logger zap.Logger
	Tracer trace.Tracer
}

func (s Serve) main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	publicKey0, err := cfg.ReadPublicKey(user.Driver)
	if err != nil {
		logger.Fatal("could not read driver public key")
	}

	publicKey1, err := cfg.ReadPublicKey(user.Passenger)
	if err != nil {
		logger.Fatal("could not read passenger public key")
	}

	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience:    cfg.DriverSalt,
			snappids.PassengerAudience: cfg.PassengerSalt,
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience:    cfg.DriverHashLength,
			snappids.PassengerAudience: cfg.PassengerHashLength,
		},
	}

	allowedAccessTypes, err := cfg.GetAllowedAccessTypes()
	if err != nil {
		logger.Fatal("error while getting allowed access types", zap.Error(err))
	}

	rest := api.API{
		Authenticator: &authenticator.Authenticator{
			PublicKeys: map[user.Issuer]*rsa.PublicKey{
				user.Driver:    publicKey0,
				user.Passenger: publicKey1,
			},
			AllowedAccessTypes: allowedAccessTypes,
			Company:            cfg.Company,
			TopicManager:       topics.NewTopicManager(cfg.Topics, hid, cfg.Company),
		},
		Tracer: tracer,
		Logger: *logger.Named("api"),
	}.ReSTServer()

	go func() {
		if err := rest.Listen(fmt.Sprintf(":%d", cfg.HTTPPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to run REST HTTP server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := rest.Shutdown(); err != nil {
		logger.Error("error happened during REST API shutdown", zap.Error(err))
	}
}

// Register serve command.
func (s Serve) Register(root *cobra.Command, cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "serve",
			Short: "serve runs the application",
			Long:  `serve will run Soteria ReST server and waits until user disrupts.`,
			Run: func(cmd *cobra.Command, args []string) {
				s.main(cfg, logger, tracer)
			},
		},
	)
}
