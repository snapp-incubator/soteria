package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/soteria/configs"
	"gitlab.snapp.ir/dispatching/soteria/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/internal/db/redis"
	"gitlab.snapp.ir/dispatching/soteria/pkg/user"
	"io/ioutil"
)

var Init = &cobra.Command{
	Use:     "init [json file]",
	Short:   "init accounts data from a json file",
	Long:    `init is used to set all accounts data from a json file `,
	Args:    cobra.ExactArgs(1),
	PreRunE: initPreRun,
	RunE:    initRun,
}

func initPreRun(cmd *cobra.Command, args []string) error {
	cfg := configs.InitConfig()

	rClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("could not init redis: %w", err)
	}

	app.GetInstance().SetAccountsService(&accounts.Service{
		Handler: redis.ModelHandler{
			Client: rClient,
		},
	})

	return nil
}

func initRun(cmd *cobra.Command, args []string) error {
	file := args[0]

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var users []user.User
	if err := json.Unmarshal(raw, &users); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	for _, u := range users {
		if err := app.GetInstance().AccountsService.Handler.Save(u); err != nil {
			return fmt.Errorf("failed to import user %s: %w", u.Username, err)
		}
	}

	return nil
}
