package cmd

import (
	"os"

	"github.com/snapp-incubator/soteria/internal/cmd/serve"
	"github.com/snapp-incubator/soteria/internal/config"
	"github.com/snapp-incubator/soteria/internal/logger"
	"github.com/snapp-incubator/soteria/internal/profiler"
	"github.com/snapp-incubator/soteria/internal/tracing"
	"github.com/spf13/cobra"
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

	profiler.Start(cfg.Profiler)

	//nolint: exhaustruct
	root := &cobra.Command{
		Use:   "soteria",
		Short: "Soteria is the authentication service.",
		Long:  `Soteria is responsible for Authentication and Authorization of every request witch send to EMQ Server.`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println("Run `soteria serve` to start serving requests")
		},
	}

	serve.Serve{
		Cfg:    cfg,
		Logger: logger.Named("serve"),
		Tracer: tracer,
	}.Register(root)

	err := root.Execute()
	if err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
