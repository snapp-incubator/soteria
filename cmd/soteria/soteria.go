package main

import (
	"log"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal/commands"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/commands/accounts"
)

func main() {
	cli := commands.Root

	cli.AddCommand(commands.Serve)

	accounts.Accounts.AddCommand(accounts.Init)
	cli.AddCommand(accounts.Accounts)
	cli.AddCommand(commands.Token)

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
