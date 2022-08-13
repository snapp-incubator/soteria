package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/soteria/internal/cmd/serve"
	"gitlab.snapp.ir/dispatching/soteria/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/internal/logger"
	"gitlab.snapp.ir/dispatching/soteria/internal/tracing"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger).Named("root")

	tracer := tracing.New(cfg.Tracer, logger.Named("tracer"))

	//nolint: exhaustruct
	root := &cobra.Command{
		Use:   "soteria",
		Short: "Soteria is the authentication service.",
		Long:  `Soteria is responsible for Authentication and Authorization of every request witch send to EMQ Server.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Run `soteria serve` to start serving requests")
		},
	}

	serve.Serve{
		Cfg:    cfg,
		Logger: *logger.Named("serve"),
		Tracer: tracer,
	}.Register(root)

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
