![Lint](https://github.com/fredbi/go-cli/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/go-cli/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/go-cli/badge.svg)](https://coveralls.io/github/fredbi/go-cli)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/go-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-cli.svg)](https://pkg.go.dev/github.com/fredbi/go-cli)

# go-cli

This repo exposes a few utilities to build command-line utilities and manage configurations on top of 
3 great librabries: `cobra`, `viper` and `pflag`. These libraries come in handy to build "12-factors" applications.

This is based on the seminal work by @casualjim. I am grateful to him for his much inspiring code.

## CLI

### Goals

The `cli` packages proposes an approach to building command-line binaries on top of the very rich `github.com/spf13/cobra` package.

There a few great libraries around to build a CLI (see below, a few that I like). `cobra` stands out as most likely the richest and most flexible,
as CLIs are entirely built programmatically.

This is great, but after some time spent building CLIs again and again, it became cleat that dealing with CLI flags and configs came with
a lot of repeatitive boiler-plate patterns.

I felt the need to take side and adopt a few opinions to build more elegant CLIs.

The goals for this lib are:
* to integrate a CLI easily with config files (see [configs](#Configuration))
* to make all config exposed through a `viper` registry
* to allow CLI flags to override this config
* to remove the boiler-plate code needed to register, then bind the flags to the config registry
* to work more easily with flags of various types, including slices or custom flag types
* to remove the need for the typical `init()` to perform all this initialization
* to remove the need to use package-level variables
* to design with testability in mind

Non-goals:
* to use struct tags: we want to stick to the programmatic approach - there are other great libraries around following the struct tags approach
* to use codegen: we want our code to be readable, not generated

Sample CLI-building code. This is an excerpt taken from [this testable example](cli/example_test.go):
```go
import (
    "fmt"

    "github.com/fredbi/go-cli/cli"
    "github.com/spf13/cobra"
)

const (
	keyLog      = "app.log.level"
	keyDry      = "run.dryRun"
)

// RootCmd builds a runnable root command
func RootCmd() *cli.Command {
	return cli.NewCommand(
		&cobra.Command{
			Use:   "example",
			Short: "examplifies a cobra command",
			Long:  "...",
			RunE:  rootRunFunc,
		},
		cli.WithFlag(nil, "dry-run", false, "Dry run",
			cli.BindFlagToConfig(keyDry),
		),
		cli.WithPersistentFlag(nil, "log-level", "info", "Controls logging verbosity",
			cli.BindFlagToConfig(keyLog),
		),
		// apply viper config to the command tree
		cli.WithConfig(cli.Config()),
	)
}

// rootRunFunc runs the root command
func rootRunFunc(c *cobra.Command, _ []string) error {
	cfg := cli.ConfigFromContext(c.Context())
	if cfg == nil {
		cli.Die("failed to retrieve config")
	}

	fmt.Println(
		"example called\n",
		fmt.Sprintf("dry-run: %t\n", cfg.GetBool(keyDry)),
		fmt.Sprintf("log level config: %s\n", cfg.GetString(keyLog)),
	)

	return nil
}

func main() {
	if err := RootCmd().Execute(); err != nil {
		cli.Die("executing: %v", err)
	}
}
```

## Configuration

The `config` package proposes an opinionated approach to dealing with config files,
and exposes a few configuration loaders that know about environments and secrets.

### Goals

The goals for this package are as follows:

* to load configuration files, using sensible defaults from the powerful `github.com/spf13/viper` package.
* to merge configurations, overloading value for a specific environment
* to deal with the specifics of secrets in config
* to help with testing the programs that consume such configuations
* to save most of the boiler plate need to deal with viper configuration settings and merging.

### Approach to configuration

The goal of this approach is to merge configuration files with environment specifics.

Applications will then be able to consume the settings from a viper configuration registry.

Supported format: YAML, JSON

Supported file extensions: "yml", "yaml", "json"

Folder structure:
```
# <- root configuration
  {base path}/config.yaml
             # <- environment-specifics
              config.d/
                       # <- configuration to merge for environment
                       {environment}/config.yaml
                       # other environment-specifics ....
                       {...}/config.yaml
```

When using default settings (these are configurable), the base path is defined by the `CONFIG_DIR` environment variable.

Secret configurations:
```
{base path}/secrets.yaml
           config.d/
                    # <- configuration to merge for environment
                    {environment}/secrets.yaml
```

### Features

* sensible defaults for minimal boiler plate
* most defaults are configurable
* extensible
* the viper object may be watched dynamically
* can merge plain config with secrets
* helper method to easily inject configurations in test programs

### Typical configuration of a kubernetes deployment

Typically, the configuration files are held in one or several `Configmap` resources, mounted by your deployed container.
Most likely, secret files will be mounted from `Secret` resources.

(... to be continued...)

### Side notes

(... to be continued...)
