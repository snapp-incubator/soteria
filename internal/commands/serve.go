package commands

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/spf13/cobra"
	snappids "gitlab.snapp.ir/dispatching/snappids/v2"
	"gitlab.snapp.ir/dispatching/soteria/configs"
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/authenticator"
	"gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"gitlab.snapp.ir/dispatching/soteria/web/grpc"
	"gitlab.snapp.ir/dispatching/soteria/web/rest/api"
	"log"
	"net"
	"os"
	"os/signal"
)

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "serve runs the application",
	Long:  `serve will run Soteria REST and gRPC server and waits until user disrupts.`,
	Run:   serveRun,
}

func serveRun(cmd *cobra.Command, args []string) {
	cfg := configs.InitConfig()
	pk, err := cfg.ReadPrivateKey(user.ThirdParty)
	if err != nil {
		log.Fatal("could not read third party private key")
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

	rClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("could not create redis client: %v", err)
	}

	app.GetInstance().SetAccountsService(&accounts.Service{
		Handler: redis.RedisModelHandler{
			Client: rClient,
		},
	})

	app.GetInstance().SetAuthenticator(&authenticator.Authenticator{
		PrivateKeys: map[string]*rsa.PrivateKey{
			user.ThirdParty: pk,
		},
		AllowedAccessTypes: cfg.GetAllowedAccessTypes(),
		ModelHandler: redis.RedisModelHandler{
			Client: rClient,
		},
		HashIDSManager:  hid,
		EMQTopicManager: snappids.NewEMQManager(hid),
	})

	rest := api.RestServer(cfg.HttpPort)

	gListen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		log.Fatal("failed to listen: %w", err)
	}

	grpcServer := grpc.GRPCServer()

	go func() {
		if err := rest.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(gListen); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	if err := rest.Shutdown(context.Background()); err != nil {
		log.Fatal("error happened during REST API shutdown: %w", err)
	}

	grpcServer.Stop()
}
