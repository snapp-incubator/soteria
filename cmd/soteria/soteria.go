package main

import (
	"gitlab.snapp.ir/dispatching/soteria/v3/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
