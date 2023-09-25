![Lint](https://github.com/fredbi/go-cli/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/go-cli/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/go-cli/badge.svg)](https://coveralls.io/github/fredbi/go-cli)
![Vulnerability Check](https://github.com/fredbi/go-cli/actions/workflows/03-govulncheck.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/go-cli)](https://goreportcard.com/report/github.com/fredbi/go-cli)

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/go-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/go-cli.svg)](https://pkg.go.dev/github.com/fredbi/go-cli)
[![license](http://img.shields.io/badge/license/License-Apache-yellow.svg)](https://raw.githubusercontent.com/fredbi/go-cli/master/LICENSE.md)

# go-cli

This repo exposes a few utilities to
(i) [build command-line utilities](#CLI) with
(ii) [flexible configurations](#Configuration)
on top of 3 great libraries:
[`github.com/spf13/cobra`](https://github.com/spf13/cobra),
[`github.com/spf13/viper`](https://github.com/spf13/viper), and
[`github.com/spf13/pflag`](https://github.com/spf13/pflag).

**TL,DR**: this is not yet another CLI-building library, but rather a mere wrapper on top of `cobra`
to use that great lib with a functional style.

## CLI

### Example for CLI

Sample CLI-building code. This example is taken from [one of the testable examples](cli/example_test.go).

Notice our main objectives here:
* no globals
* inline flag registration & binding
* access to settings using viper only

```go
package main

import (
    "fmt"

    "github.com/fredbi/go-cli/cli"
    "github.com/spf13/cobra"
)

const (
    // viper config keys
	keyLog      = "app.log.level"
	keyDry      = "run.dryRun"
)

func main() {
    // no global vars, no init() ...
	if err := RootCmd().Execute(); err != nil {
		cli.Die("executing: %v", err)
	}
}

// RootCmd builds a runnable root command
func RootCmd() *cli.Command {
	return cli.NewCommand(
        // your usual cobra command, wrapped as a function
		&cobra.Command{
			Use:   "example",
			Short: "examplifies a cobra command",
			Long:  "...",
			RunE:  rootRunFunc,
		},
        // flag bindings
        // {flag name}, {the flag type is inferred from the default value}, {flag help description}
		cli.WithFlag("dry-run", false, "Dry run",
			cli.BindFlagToConfig(keyDry), // flag bindings to a viper config
		),
        // a flag inherited by subcommands
		cli.WithPersistentFlag("log-level", "info", "Controls logging verbosity",
			cli.BindFlagToConfig(keyLog),
		),
        // apply viper config to the command tree
        // command binding to a viper config -> config will be available from context
		cli.WithConfig(cli.Config()),
	)
}

// rootRunFunc runs the root command
func rootRunFunc(c *cobra.Command, _ []string) error {
    // retrieve injected dependencies, create new empty viper registry if unresolved
	cfg := injectable.ConfigFromContext(c.Context(), viper.New)

	fmt.Println(
		"example called\n",
		fmt.Sprintf("dry-run: %t\n", cfg.GetBool(keyDry)),
		fmt.Sprintf("log level config: %s\n", cfg.GetString(keyLog)),
	)

	return nil
}
```

### Goals

The `cli` packages proposes an opinionated approach to building command-line binaries on top of `github.com/spf13/cobra`.

> There are a few great existing libraries around to build a CLI.
> I believe that `cobra` stands out as the richest and most flexible,
> as CLIs are entirely built programmatically.

`cobra` is great, but building CLIs again and again, I came to identify a few repetitive boiler-plate patterns.

So this module reflects my opinions about how to build more elegant CLIs, wich abide by [12-factor](https://12factor.net)
out-of-the-box, with more expressive code and less low-level tinkering.

Feedback is always welcome, as opinions may evolve over time...
Feel free to post issues to leave your comments and/or proposals.

[More detailed design goals](docs/goals.md#CLI)

## Configuration

The `config` package proposes an opinionated approach to dealing with config files on top of `github.com/spf13/viper`.

It exposes configuration loaders which know about the deployment context
(e.g a deployment environment such as `dev`, `production`) and secrets.

Although developped primarily to serve a CLI, this package may be used independently.

### Example: loading a config

Other examples are available [here](config/examples_test.go).

```go
import (
	"fmt"
	"log"

	"github.com/fredbi/go-cli/config"
)

func ExampleLoad() {
    ...

	// load and merge configuration files for environment "dev"
	cfg, err := config.Load("dev", config.WithMute(true))
	if err != nil {
		log.Fatalf("loading config: %w", err)

		return
	}

	fmt.Println(cfg.AllSettings())
}
```

### Goals

This describes my approach to configuration. We want to:

1. retrieve a config organized as a hierarchy of settings, e.g. a YAML document
2. merge configuration files with environment-specific settings
3. merge configuration files with secrets, usually these are environment-specific
4. clearly isolate and merge default settings
5. applications to be able to consume the settings from a single viper configuration registry

In addition,

* we want the hierarchy to be agnostic to the environment context
* most of the time, we don't want env-specific sections to propagate to the app level
  (e.g. in the style of `.ini` sections)

> In our code, we should never check for a dev or prod specific section of the configuration.


Supported format: YAML, JSON

Supported file extensions: "yml", "yaml", "json"

See other [examples](.config/examples_test.go)

[More detailed design goals](docs/goals.md#Configuration)

### Folders structure for configurations

By default we have:
```
# <- root configuration
{base path}/config.yaml
            # <- environment-specifics folder
            config.d/
                     # <- extra configuration to merge
                     config.yaml
                     # <- possibly with a modified name: config.*.yaml
                     config.default.yaml
                     # <- configuration to merge for environment
                     {environment}/config.yaml
                     # other environment-specifics ....
                     {...}/config.yaml
```

[Here is an example](./config/examples)

When using default settings for this module (these are configurable),
the base path is defined by the `CONFIG_DIR` environment variable.

Secret configurations:
```
{base path}/secrets.yaml
            config.d/
                     # <- secrets to merge
                     secrets.yaml
                     # <- configuration to merge for environment
                     {environment}/secrets.yaml
```

### [Typical configuration for a Kubernetes deployment](docs/k8s.md)

### Side notes

#### TODOs

* [CLI todo list](cli/TODO.md)
* [config todo list](config/TODO.md)

#### Dealing with secrets locally

TODO(fredbi)

### Credits

The config part is largely based on some seminal past work by [@casualjim](https://github.com/casualjim/).
I am grateful to him for his much inspiring code.

The version from runtime piece of code is largely inspired by the wonderful work from the
[golangci](https://github.com/golangci/golangci-lint) community.
