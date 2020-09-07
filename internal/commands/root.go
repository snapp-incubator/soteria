package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "soteria",
	Short: "Soteria is the authentication service.",
	Long:  `Soteria is responsible for Authentication and Authorization of every request witch send to EMQ Server.`,
	Run:   rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	fmt.Println("Run `soteria serve` to start serving requests")
}
