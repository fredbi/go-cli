//nolint:forbidigo
package cli_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fredbi/go-cli/cli"
	"github.com/fredbi/go-cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// globalFlags captures CLI flags.
//
// In this example, we choose to control over where the flag values are stored.
//
// This is not needed if all configuration is bound to viper.
var globalFlags cliFlags

func init() {
	// set some config options for testing.
	here, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current working dir: %v:", err)

		return
	}

	// we want to look for config files in the "fixtures" folder.
	fixtures := filepath.Join(here, "fixtures")
	cli.SetConfigOptions(
		config.WithMute(true),
		config.WithWatch(false),
		config.WithBasePath(fixtures),
	)
}

type (
	// TODO: could use struct tags for default and config key
	cliFlags struct {
		DryRun   bool
		URL      string
		Parallel int
		User     string
		LogLevel string

		Child childFlags // flags for the child command
	}

	childFlags struct {
		Workers int
	}
)

// Default values for flags.
func (f cliFlags) Defaults() cliFlags {
	return cliFlags{
		URL:      "https://www.example.com",
		Parallel: 2,
		LogLevel: "info",
		Child: childFlags{
			Workers: 5,
		},
	}
}

// applyDefaults set default values for the config. It is consistent with flag defaults.
func (f cliFlags) applyDefaults(cfg *viper.Viper) {
	cfg.SetDefault(keyURL, globalFlags.Defaults().URL)
	cfg.SetDefault(keyParallel, globalFlags.Defaults().Parallel)
	cfg.SetDefault(keyLog, globalFlags.Defaults().LogLevel)
	cfg.SetDefault(keyWorkers, globalFlags.Defaults().Child.Workers)
}

// nolint: dupl
// RootCmdWithGlobal illustrates the scaffolding of a command tree with explicit storage
// of the CLI flags and default values.
//
// Although not my preferred style, we can still do it this way and do something efficient.
func RootCmdWithGlobal() *cli.Command {
	return cli.NewCommand(
		&cobra.Command{
			Use:   "example",
			Short: "examplifies a cobra command with global flags",
			Long:  "...",
			RunE:  rootGlobalRunFunc,
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
					RunE:  childGlobalRunFunc,
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

// root command execution
func rootGlobalRunFunc(c *cobra.Command, args []string) error {
	err := rootRunFunc(c, args)

	fmt.Println(
		"global flags values evaluated by root\n",
		fmt.Sprintf("%#v", globalFlags),
	)

	return err
}

// child command execution
func childGlobalRunFunc(c *cobra.Command, args []string) error {
	err := childRunFunc(c, args)

	fmt.Println(
		"global flags values evaluated by child\n",
		fmt.Sprintf("%#v", globalFlags),
	)

	return err
}

func Example_globalRootCmd() {
	rootCmd := RootCmdWithGlobal()
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
func Example_globalChildCmd() {
	rootCmd := RootCmdWithGlobal()
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
