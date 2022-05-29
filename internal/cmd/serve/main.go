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

func (s Serve) main() {
	publicKey0, err := s.Cfg.ReadPublicKey(user.Driver)
	if err != nil {
		s.Logger.Fatal("could not read driver public key")
	}

	publicKey1, err := s.Cfg.ReadPublicKey(user.Passenger)
	if err != nil {
		s.Logger.Fatal("could not read passenger public key")
	}

	hid := &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience:    s.Cfg.DriverSalt,
			snappids.PassengerAudience: s.Cfg.PassengerSalt,
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience:    s.Cfg.DriverHashLength,
			snappids.PassengerAudience: s.Cfg.PassengerHashLength,
		},
	}

	allowedAccessTypes, err := s.Cfg.GetAllowedAccessTypes()
	if err != nil {
		s.Logger.Fatal("error while getting allowed access types", zap.Error(err))
	}

	rest := api.API{
		Authenticator: &authenticator.Authenticator{
			PublicKeys: map[user.Issuer]*rsa.PublicKey{
				user.Driver:    publicKey0,
				user.Passenger: publicKey1,
			},
			AllowedAccessTypes: allowedAccessTypes,
			Company:            s.Cfg.Company,
			TopicManager:       topics.NewTopicManager(s.Cfg.Topics, hid, s.Cfg.Company),
		},
		Tracer: s.Tracer,
		Logger: *s.Logger.Named("api"),
	}.ReSTServer()

	go func() {
		if err := rest.Listen(fmt.Sprintf(":%d", s.Cfg.HTTPPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Fatal("failed to run REST HTTP server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := rest.Shutdown(); err != nil {
		s.Logger.Error("error happened during REST API shutdown", zap.Error(err))
	}
}

// Register serve command.
func (s Serve) Register(root *cobra.Command) {
	root.AddCommand(
		// nolint: exhaustruct
		&cobra.Command{
			Use:   "serve",
			Short: "serve runs the application",
			Long:  `serve will run Soteria ReST server and waits until user disrupts.`,
			Run: func(cmd *cobra.Command, args []string) {
				s.main()
			},
		},
	)
}
