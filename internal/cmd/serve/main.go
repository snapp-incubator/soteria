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
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/api"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/metrics"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/tracer"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger) {
	publicKey0, err := cfg.ReadPublicKey(user.Driver)
	if err != nil {
		zap.L().Fatal("could not read driver public key")
	}

	publicKey1, err := cfg.ReadPublicKey(user.Passenger)
	if err != nil {
		zap.L().Fatal("could not read passenger public key")
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

	trc, cl, err := tracer.New(cfg.Tracer)
	if err != nil {
		zap.L().Fatal("could not create tracer", zap.Error(err))
	}

	app.GetInstance().SetTracer(trc, cl)

	allowedAccessTypes, err := cfg.GetAllowedAccessTypes()
	if err != nil {
		zap.L().Fatal("error while getting allowed access types", zap.Error(err))
	}

	app.GetInstance().SetAuthenticator(&authenticator.Authenticator{
		PublicKeys: map[user.Issuer]*rsa.PublicKey{
			user.Driver:    publicKey0,
			user.Passenger: publicKey1,
		},
		AllowedAccessTypes: allowedAccessTypes,
		ModelHandler:       db.NewInternal(cfg.Users),
		HashIDSManager:     hid,
		EMQTopicManager:    snappids.NewEMQManagerWithCompany(hid, cfg.Company),
		Company:            cfg.Company,
	})

	m := metrics.NewMetrics()
	app.GetInstance().SetMetrics(&m.Handler)

	rest := api.ReSTServer()

	go func() {
		if err := rest.Listen(fmt.Sprintf(":%d", cfg.HTTPPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("failed to run REST HTTP server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := rest.Shutdown(); err != nil {
		zap.L().Error("error happened during REST API shutdown", zap.Error(err))
	}

	if err := app.GetInstance().TracerCloser.Close(); err != nil {
		zap.L().Error("error happened while closing tracer", zap.Error(err))
	}
}

// Register serve command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "serve",
			Short: "serve runs the application",
			Long:  `serve will run Soteria ReST server and waits until user disrupts.`,
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger)
			},
		},
	)
}
