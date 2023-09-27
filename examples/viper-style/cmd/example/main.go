package main

import (
	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/examples/viper-style/cmd/example/commands"
)

func main() {
	cli.MustOrDie("executing example",
		commands.Root.Execute(),
	)
}
