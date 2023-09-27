package cli_test

import (
	"fmt"
	"log"

	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Example_defaults() {
	if err := RootCmdWithDefaults().
		ExecuteWithArgs("defaults"); err != nil {
		log.Fatal("executing:", err)
	}

	// Output:
	// map[app:map[log:map[level:info]] nonflagkey:10.345]
}

func RootCmdWithDefaults() *cli.Command {
	return cli.NewCommand(
		&cobra.Command{
			Use:   "defaults",
			Short: "examplifies a cobra command with defaults settings",
			Long:  "this command shows how config defaults can be set from flags, and supplemented with non-flag defaults",
			RunE:  rootWithDefaultsRunFunc,
		},
		cli.WithFlag("log-level", "info", "Controls logging verbosity",
			cli.FlagIsPersistent(),
			cli.BindFlagToConfig(keyLog),
		),
		cli.WithConfig(cli.Config(func(v *viper.Viper) {
			v.SetDefault("nonFlagKey", 10.345)
		})),
		cli.WithAutoVersion(),
	)
}

func rootWithDefaultsRunFunc(cmd *cobra.Command, _ []string) error {
	cfg := injectable.ConfigFromContext(cmd.Context(), viper.New)

	fmt.Println(cfg.AllSettings())

	return nil
}
