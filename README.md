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
on top of these 3 great libraries:
`github.com/spf13/cobra`, `github.com/spf13/viper` and `github.com/spf13/pflag`.

The config part is based on some seminal past work by [@casualjim](https://github.com/casualjim/).
I am grateful to him for his much inspiring code.

**TL,DR**: this is not yet another CLI-building library, but rather a mere wrapper on top of `cobra`
to use that great lib with a better style (IMHO).

## CLI

### Example for CLI

Sample CLI-building code. This is an excerpt taken from [this testable example](cli/example_test.go):

```go
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
		cli.WithFlag("dry-run", false, "Dry run", // {flag name}, {flag type inferred from the default value}, {flag help description}
			cli.BindFlagToConfig(keyDry), // flag bindings to a viper config
		),
		cli.WithPersistentFlag("log-level", "info", "Controls logging verbosity",
			cli.BindFlagToConfig(keyLog),
		),
		// apply viper config to the command tree
		cli.WithConfig(cli.Config()), // command binding to a viper config -> will be available from context
	)
}

// rootRunFunc runs the root command
func rootRunFunc(c *cobra.Command, _ []string) error {
	cfg := cli.ConfigFromContext(c.Context()) // retrieve the config
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
    // no global vars, no init() ...
	if err := RootCmd().Execute(); err != nil {
		cli.Die("executing: %v", err)
	}
}
```
### Goals

The `cli` packages proposes an approach to building command-line binaries on top of `github.com/spf13/cobra`.

> There are a few great existing libraries around to build a CLI (see below, a few that I like).
> `cobra` stands out as the richest and most flexible,
> as CLIs are entirely built programmatically.

`cobra` is great, but building CLIs again and again, I came to identify a few repetitive boiler-plate patterns.

So this module reflects my opinions about how to build more elegant CLIs, with more expressive code and less tinkering.

Feedback is always welcome, as opinions may evolve over time...
Feel free to post issues to leave your comments and/or proposals.

#### Desirable features

* a typical CLI should interact easily with config files (see [Configuration](#Configuration)), but not _always_
  * expose all config through a `viper` registry
  * leave developers a free-hand on all the knobs and features proposed by `cobra`
* it should be easier to interact with command line flags of various types
  * simple, declarative registration and binding of flags to config
  * should abstract away the tedious and error prone steps for the registration of flags, binding & defaults
  * allow CLI flags to override this config (12-factors)
  * includes slices, maps and custom flag types (delegated to `github.com/fredbi/gflag`)
* it should be easier to declare defaults (for flags, for config)
* it should be easier to inject external dependencies into the commands tree (config, logger, etc)

#### Code style goals

* adopt a functional style (some would say DSL-like), with builder functions for CLI components 
* the tree-like structure of commands should appear visually and obviously in the source code
* remove the need for the typical `init()` to perform all this initialization
* remove the need to use package-level variables
* remove the boiler-plate code needed to register, then bind the flags to the config registry
* remove the cognitive burden of remembering all the `GetString()`, `GetBool()` etc methods: `go` now has generics for that
* favor the use of generics exposed by `github.com/fredbi/gflags`, but don't require it
* design with testability in mind: CLI's should be testable with reasonable code coverage

#### Non-goals

* don't use struct tags: we want to stick to the programmatic approach - there are other great libraries around following the struct tags approach
* don't use codegen: we want our code to be readable, not generated


## Configuration

The `config` package proposes an opinionated approach to dealing with config files on top of `github.com/spf13/viper`.

It exposes configuration loaders that know about the context
(e.g a deployment environment such as `dev`, `production`) and secrets.

Although developped primarily to serve a CLI, this package may be used independently.

### Example: loading a config

```go
import (
	"fmt"
	"log"
	"os"

	"github.com/fredbi/go-cli/config"
)

func ExampleLoad() {
	here := mustCwd()
	os.Setenv("CONFIG_DIR", filepath.Join(here, "examples"))

	// load and merge configuration files for environment "dev"
	cfg, err := config.Load("dev", config.WithMute(true))
	if err != nil {
		err = fmt.Errorf("loading config: %w", err)
		log.Fatal(err)

		return
	}

	fmt.Println(cfg.AllSettings())
}

func mustCwd() string {
	here, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("get current working dir: %w", err)
		log.Fatal(err)
	}

	return here
}
```

### Goals

#### Desirable features

* load configuration files, using sensible defaults from the powerful `github.com/spf13/viper` package.
* merge configurations, overloading value for a specific environment
* deal with the specifics of merge secrets in config
* help with testing the programs that consume configurations
* leverages all the 12-factor app stuff from `viper`
* leave developers a free-hand if they want to use all the knobs and features proposed by `viper`
  (e.g. dynamic watch, remote config, etc)
* defaults are configurable

#### Code-style goals

* less boiler plate to deal with `viper` configuration settings and merging

#### Non-goals

* Avoid too much of automagically resolving things

  > As much as a like what's available for pythonists with [Dynaconf](https://www.dynaconf.com/),
  > I found myself spending too much time reading their doc to understand their default settings
  > and figure out whetheir to adopt the default or override.
  >
  > This is a pitfall that is very difficult to avoid, and only experience and feedback will tell.


* Don't want to support older config formats such as `.ini`, `.toml` etc

  > While perfectly doable, I prefer at the moment to focus on having less things to describe and
  > document. So I believe that YAML and JSON are good enough.

* At the moment, no particular goal is set to support secrets via APIs (e.g. Hashicorp's Vault, Azure Vault...)

  > Let's wait for a bit.
  > At the moment, I am assuming secrets are just plain files (e.g. Kubernetes secret)


### Approach to configuration

We want to:
1. retrieve a config organized as a hierarchy of settings, e.g. a YAML document
2. merge configuration files with environment-specific settings
3. merge configuration files with secrets, usually these are environment-specific
4. clearly isolate and merge default settings

In addition,

* we want the hierarchy to be agnostic to the environment context
* most of the time, we don't want env-specific sections to propagate to the app level
  (e.g. in the style of `.ini` sections)

> In our code, we should never check for a dev or prod specific section of the configuration.

Applications are able to consume the settings from a single viper configuration registry.

Supported format: YAML, JSON

Supported file extensions: "yml", "yaml", "json"

See other [examples](.config/examples_test.go)

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

### Typical configuration for a Kubernetes deployment

Typically, the configuration files are held in one or several `Configmap` resources, mounted by your deployed container.

Secret files can be mounted from `Secret` resources in the container, accessible as plain files.

Alternatively, Kubernetes may expose secrets as environment variables: `viper` takes care of loading them in the registry.

Normally, we don't want to expose secrets via CLI flags.

Example (e.g. volumes & container section of a k8s PodTemplateSpec):
```yaml
volumes:
  - name: config
    configMap:          # <- expose config file from ConfigMap resource to the pod's containers
      name: 'app-config'
  - name: secret-config # <- expose secrets file from Secret as file resource to the pod's containers
    secret:
      secret_name: 'app-secret-config'

containers:
  - name: app-container
    ...
    env:
      - name: CONFIG_DIR
        value: '/etc/app'
      - name: SECRET_URL
        valueFrom:
          secretKeyRef: # <- expose config value as an environment variable to the container
          name: 'app-secret-url'
          key: secretUrl

    volumeMounts:
      - mountPath: '/etc/app' # <- mount config file(s) as /etc/app/{key(s)} file(s)
        name: config
      - mountPath: '/etc/app/config.d'
```

### Side notes

#### Dealing with secrets locally

TODO(fredbi)
