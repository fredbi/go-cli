//nolint:forbidigo
package cli_test

import (
	"fmt"
	"log"

	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// config key mappings
const (
	/*
		YAML config accepted:

		app:
		  url: xyz
		  parallel: 99
		  user: abc
		  child:
		    workers: 5
		log:
		  level: info
	*/
	keyLog      = "app.log.level"
	keyURL      = "app.url"
	keyParallel = "app.parallel"
	keyUser     = "app.user"
	keyWorkers  = "app.child.workers"
	keyDry      = "run.dryRun"
)

// root command execution
func rootRunFunc(c *cobra.Command, _ []string) error {
	// retrieve injected dependencies, create new empty viper registry if unresolved
	cfg := injectable.ConfigFromContext(c.Context(), viper.New)

	fmt.Println(
		"example called\n",
		fmt.Sprintf("URL config: %s\n", cfg.GetString(keyURL)),
		fmt.Sprintf("log level config: %s\n", cfg.GetString(keyLog)),
		fmt.Sprintf("parallel config: %d\n", cfg.GetInt(keyParallel)),
		fmt.Sprintf("user config: %s\n", cfg.GetString(keyUser)),
	)

	fmt.Println(
		"global flags values evaluated by root\n",
		fmt.Sprintf("%#v", globalFlags),
	)

	return nil
}

// child command execution
func childRunFunc(c *cobra.Command, _ []string) error {
	cfg := injectable.ConfigFromContext(c.Context(), viper.New)

	fmt.Println(
		"subcommand called\n",
		fmt.Sprintf("URL config: %s\n", cfg.GetString(keyURL)),
		fmt.Sprintf("parallel config: %d\n", cfg.GetInt(keyParallel)),
		fmt.Sprintf("user config: %s\n", cfg.GetString(keyUser)),
		fmt.Sprintf("workers config: %s\n", cfg.GetString(keyWorkers)),
	)

	fmt.Println(
		"global flags values evaluated by child\n",
		fmt.Sprintf("%#v", globalFlags),
	)

	return nil
}

func emptyRunFunc(c *cobra.Command, _ []string) error {
	cfg := injectable.ConfigFromContext(c.Context(), viper.New)

	fmt.Println("command called:", c.Name())
	fmt.Println("injected config:", cfg.AllSettings())

	return nil
}

// nolint: dupl
// RootCmd illustrates the scaffolding of a command tree without any explicit storage of the CLI flags and default values.
//
// TODO: binding a flag to the config should trigger a SetDefault in the config. Problem: should be done before loading,
// HENCE config loader should be lazy.
func RootCmd() *cli.Command {

	return cli.NewCommand(
		&cobra.Command{
			Use:   "example",
			Short: "examplifies a cobra command",
			Long:  "...",
			RunE:  rootRunFunc,
		},
		cli.WithFlagVar(&globalFlags.DryRun, "dry-run", globalFlags.Defaults().DryRun, "Dry run",
			cli.BindFlagToConfig(keyDry),
		),
		cli.WithFlagVar(&globalFlags.LogLevel, "log-level", globalFlags.Defaults().LogLevel, "Controls logging verbosity",
			cli.FlagIsPersistent(),
			cli.BindFlagToConfig(keyLog),
		),
		cli.WithFlagVar(&globalFlags.URL, "url", globalFlags.Defaults().URL, "The URL to connect to",
			cli.FlagIsPersistent(),
			cli.BindFlagToConfig(keyURL),
		),
		cli.WithFlagVarP(&globalFlags.Parallel, "parallel", "p", globalFlags.Defaults().Parallel, "Degree of parallelism",
			cli.FlagIsPersistent(),
			cli.BindFlagToConfig(keyParallel),
		),
		// example with RegisterFunc, useful for maximum flexibility.
		cli.WithFlagFunc(func(flags *pflag.FlagSet) string {
			const userFlag = "user"
			flags.StringVar(&globalFlags.User, userFlag, globalFlags.Defaults().User, "Originating user")
			return userFlag
		},
			cli.FlagIsPersistent(),
			cli.FlagIsRequired(),
			cli.BindFlagToConfig(keyUser),
		),
		cli.WithSubCommands(
			cli.NewCommand(
				&cobra.Command{
					Use:   "child",
					Short: "sub-command example",
					Long:  "...",
					RunE:  childRunFunc,
				},
				cli.WithFlagVar(&globalFlags.Child.Workers, "workers", globalFlags.Defaults().Child.Workers, "Number of workers threads",
					cli.FlagIsRequired(),
					cli.BindFlagToConfig(keyWorkers),
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
			cli.NewCommand(
				&cobra.Command{
					Use:   "version",
					Short: "another sub-command example",
					Long:  "...",
					RunE:  emptyRunFunc,
				},
			),
		),
		// apply config to the command tree
		cli.WithConfig(cli.Config(globalFlags.applyDefaults)),
	)
}

func Example_help() {
	rootCmd := RootCmd()
	rootCmd.SetArgs([]string{
		"--help",
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("executing:", err)
	}

	fmt.Println("done")

	// Output:
	// 	...
	//
	// Usage:
	//   example [flags]
	//   example [command]
	//
	// Available Commands:
	//   child       sub-command example
	//   completion  Generate the autocompletion script for the specified shell
	//   help        Help about any command
	//   version     another sub-command example
	//
	// Flags:
	//       --dry-run            Dry run
	//   -h, --help               help for example
	//       --log-level string   Controls logging verbosity (default "info")
	//   -p, --parallel int       Degree of parallelism (default 2)
	//       --url string         The URL to connect to (default "https://www.example.com")
	//       --user string        Originating user
	//
	// Use "example [command] --help" for more information about a command.
	// done
}

// Example_rootCmd runs the root command.
func Example_rootCmd() {
	rootCmd := RootCmd()
	rootCmd.SetArgs([]string{
		"--dry-run",
		"--log-level",
		"debug",
		"--parallel",
		"15",
		"--user",
		"fred",
	},
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("executing:", err)
	}

	fmt.Println("done")

	// Output:
	// example called
	//  URL config: https://www.example.com
	//  log level config: debug
	//  parallel config: 15
	//  user config: fred
	//
	// global flags values evaluated by root
	//  cli_test.cliFlags{DryRun:true, URL:"https://www.example.com", Parallel:15, User:"fred", LogLevel:"debug", Child:cli_test.childFlags{Workers:5}}
	// done
}

// Example_childCmd runs a child command.
func Example_childCmd() {
	rootCmd := RootCmd()
	if err := rootCmd.ExecuteWithArgs(
		"child",
		"--parallel",
		"20",
		"--url",
		"https://www.zorg.com",
		"--user",
		"zorg",
		"--workers",
		"12",
	); err != nil {
		log.Fatal("executing:", err)
	}

	fmt.Println("done")

	// Output:
	// subcommand called
	//  URL config: https://www.zorg.com
	//  parallel config: 20
	//  user config: zorg
	//  workers config: 12
	//
	// global flags values evaluated by child
	//  cli_test.cliFlags{DryRun:false, URL:"https://www.zorg.com", Parallel:20, User:"zorg", LogLevel:"info", Child:cli_test.childFlags{Workers:12}}
	// done

}

// Example_printCmd prints out the tree structure of the root command.
func Example_printCmd() {
	rootCmd := RootCmd()

	fmt.Println(rootCmd.String())

	// Output:
	// example
	//  child
	//   grandchild
	//  version
}

// TODO: interact with RegisterFlags() style of registration in-bulk (viper.BindPflags??)
