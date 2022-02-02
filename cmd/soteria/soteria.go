package main

import (
	"gitlab.snapp.ir/dispatching/soteria/internal/cmd"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
