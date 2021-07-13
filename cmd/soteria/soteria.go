package main

import (
	"log"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal/commands"
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/commands/accounts"
)

var cli = commands.Root

func init() {
	cli.AddCommand(commands.Serve)

	accounts.Accounts.AddCommand(accounts.Init)
	cli.AddCommand(accounts.Accounts)
	cli.AddCommand(commands.Token)
}

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
