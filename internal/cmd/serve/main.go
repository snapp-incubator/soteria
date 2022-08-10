package serve

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/internal/api"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/topics"
	"gitlab.snapp.ir/dispatching/soteria/pkg/acl"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Serve struct {
	Cfg    config.Config
	Logger zap.Logger
	Tracer trace.Tracer
}

func (s Serve) main() {
	rest := api.API{
		Authenticators: s.Authenticators(),
		Tracer:         s.Tracer,
		Logger:         *s.Logger.Named("api"),
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

func (s Serve) Authenticators() map[string]*authenticator.Authenticator {
	all := make(map[string]*authenticator.Authenticator)

	for _, vendor := range s.Cfg.Vendors {
		publicKeys := s.PublicKeys(vendor.Keys)
		hid := HIDManager(vendor.DriverSalt, vendor.DriverHashLength, vendor.PassengerSalt, vendor.PassengerHashLength)
		allowedAccessTypes := s.GetAllowedAccessTypes(vendor.AllowedAccessTypes)

		auth := &authenticator.Authenticator{
			PublicKeys:         publicKeys,
			AllowedAccessTypes: allowedAccessTypes,
			Company:            vendor.Company,
			TopicManager:       topics.NewTopicManager(vendor.Topics, hid, vendor.Company),
		}

		all[vendor.Company] = auth
	}

	return all
}

func HIDManager(
	driverSalt string,
	driverHashLength int,
	passengerSalt string,
	passengerHashLength int,
) *snappids.HashIDSManager {
	return &snappids.HashIDSManager{
		Salts: map[snappids.Audience]string{
			snappids.DriverAudience:    driverSalt,
			snappids.PassengerAudience: passengerSalt,
		},
		Lengths: map[snappids.Audience]int{
			snappids.DriverAudience:    driverHashLength,
			snappids.PassengerAudience: passengerHashLength,
		},
	}
}

func (s Serve) PublicKeys(keys map[string]string) map[string]*rsa.PublicKey {
	rsaKeys := make(map[string]*rsa.PublicKey)

	for iss, publicKey := range keys {
		rsaKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			s.Logger.Fatal("could not read public key", zap.String("issuer", iss))
		}

		rsaKeys[iss] = rsaKey
	}

	return rsaKeys
}

// GetAllowedAccessTypes will return all allowed access types in Soteria.
func (s Serve) GetAllowedAccessTypes(accessTypes []string) []acl.AccessType {
	allowedAccessTypes := make([]acl.AccessType, 0, len(accessTypes))

	for _, a := range accessTypes {
		at, err := toUserAccessType(a)
		if err != nil {
			err = fmt.Errorf("could not convert %s: %w", at, err)
			s.Logger.Fatal("error while getting allowed access types", zap.Error(err))
		}

		allowedAccessTypes = append(allowedAccessTypes, at)
	}

	return allowedAccessTypes
}

// toUserAccessType will convert string access type to it's own type.
// nolint: goerr113
func toUserAccessType(access string) (acl.AccessType, error) {
	switch access {
	case "pub":
		return acl.Pub, nil
	case "sub":
		return acl.Sub, nil
	case "pubsub":
		return acl.PubSub, nil
	}

	return "", fmt.Errorf("%v is a invalid acces type", access)
}
