package accounts

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Accounts = &cobra.Command{
	Use:   "accounts",
	Short: "Soteria account management",
	Long:  `accounts is used to manage all operations related to accounts`,
	Run:   accountsRun,
}

func accountsRun(cmd *cobra.Command, args []string) {
	fmt.Println("run `soteria accounts --help` to see available options")
}
