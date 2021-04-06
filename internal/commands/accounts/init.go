package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/accounts"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/app"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/config"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/db/redis"
	"gitlab.snapp.ir/dispatching/soteria/v3/pkg/user"
	"golang.org/x/crypto/bcrypt"
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
	cfg := config.InitConfig()

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
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate hash for %s, %w", u.Username, err)
		}
		u.Password = string(hash)

		if err := app.GetInstance().AccountsService.Handler.Save(cmd.Context(), u); err != nil {
			return fmt.Errorf("failed to import user %s: %w", u.Username, err)
		}
	}

	return nil
}
