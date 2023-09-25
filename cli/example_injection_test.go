package cli_test

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
)

func RootCmdWithDep() *cli.Command {
	return cli.NewCommand(
		&cobra.Command{
			Use:   "inject",
			Short: "examplifies a cobra command with dependency injection",
			Long:  "in this example, we inject a logger in the context of the command (from log/slog)",
			RunE:  rootWithDepRunFunc,
		},
		cli.WithFlag("log-level", "info", "Controls logging verbosity",
			cli.FlagIsPersistent(),
		),
		cli.WithInjectables( // injectables know how to be retrieved from context
			injectable.NewSLogger(slog.New(slog.NewJSONHandler(os.Stderr, nil))),
		),
	)
}

func rootWithDepRunFunc(cmd *cobra.Command, _ []string) error {
	// retrieve injected dependencies from the context of the command
	ctx := cmd.Context()
	logger := injectable.SLoggerFromContext(ctx).With("command", "root")
	if logger == nil {
		return errors.New("no logger provided")
	}

	fmt.Println("print log on stderr")
	logger.Info("structured log entry")

	return nil
}

func ExampleWithInjectables() {
	rootCmd := RootCmdWithDep()
	rootCmd.SetArgs([]string{
		"inject",
		"log-level",
		"info",
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("executing:", err)
	}

	// Output:
	// print log on stderr
}
