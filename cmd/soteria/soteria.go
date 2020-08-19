package main

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/commands"
	"log"
)

var cli = commands.Root

func init() {
	cli.AddCommand(commands.Serve)
}

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
