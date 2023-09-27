package commands

import (
	"fmt"

	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/cli/injectable"
	keys "github.com/fredbi/go-cli/examples/viper-style/pkg/config-keys"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Root ...
func Root() *cli.Command {
	return cli.NewCommand(
		&cobra.Command{
			Use:   "example",
			Short: "examplifies a cobra command",
			Long:  "...",
			RunE:  rootRunFunc,
		},
		cli.WithFlag("dry-run", false, "Dry run"),
		cli.WithBindFlagsToConfig(map[string]string{
			"dry-run": keys.KeyDry,
		}),
		// persistent flags
		cli.WithFlag("log-level", "info", "Controls logging verbosity",
			cli.FlagIsPersistent(),
		),
		cli.WithFlag("url", "https://www.example.com", "The URL to connect to",
			cli.FlagIsPersistent(),
		),
		cli.WithFlagP("parallel", "p", 2, "Degree of parallelism",
			cli.FlagIsPersistent(),
		),
		cli.WithFlagFunc(func(flags *pflag.FlagSet) string {
			const userFlag = "user"
			flags.String(userFlag, "", "Originating user")

			return userFlag
		},
			cli.FlagIsPersistent(),
			cli.FlagIsRequired(),
		),
		cli.WithBindPersistentFlagsToConfig(map[string]string{
			"log-level": keys.KeyLog,
			"url":       keys.KeyURL,
			"parallel":  keys.KeyParallel,
			"user":      keys.KeyUser,
		}),
		cli.WithSubCommands(
			cli.NewCommand(
				&cobra.Command{
					Use:   "child",
					Short: "sub-command example",
					Long:  "This sub-command inherits flags from the root command and declares a local flag --workers.",
					RunE:  childRunFunc,
				},
				cli.WithFlag("workers", 5, "Number of workers threads",
					cli.FlagIsRequired(),
					cli.BindFlagToConfig(keys.KeyWorkers),
				),
				cli.WithSubCommands(
					cli.NewCommand(
						&cobra.Command{
							Use:   "grandchild",
							Short: "sub-sub-command example",
							Long:  "...",
							RunE:  emptyRunFunc,
						},
					),
				),
			),
		),
		cli.WithConfig(cli.Config()),
		cli.WithAutoVersion(),
	)
}

// root command execution
func rootRunFunc(c *cobra.Command, _ []string) error {
	// retrieve injected dependencies, create new empty viper registry if unresolved
	_ = injectable.ConfigFromContext(c.Context(), viper.New)

	return nil
}

// child command execution
func childRunFunc(c *cobra.Command, _ []string) error {
	_ = injectable.ConfigFromContext(c.Context(), viper.New)

	return nil
}

// noop command execution
func emptyRunFunc(c *cobra.Command, _ []string) error {
	cfg := injectable.ConfigFromContext(c.Context(), viper.New)

	fmt.Println("command called:", c.Name())
	fmt.Println("injected config:", cfg.AllSettings())

	return nil
}
