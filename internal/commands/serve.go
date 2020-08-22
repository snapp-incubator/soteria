package commands

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/soteria/configs"
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
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

	rClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("could not create redis client: %v", err)
	}

	app.GetInstance().SetAccountsService(&accounts.Service{
		Handler: redis.RedisModelHandler{
			Client: rClient,
		},
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
