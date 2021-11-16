package main

import (
	"log"

	"gitlab.snapp.ir/dispatching/soteria/v3/internal/commands"
)

func main() {
	cli := commands.Root

	cli.AddCommand(commands.Serve)

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
