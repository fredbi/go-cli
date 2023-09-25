# Design goals

## CLI

### Desirable features

1. a typical CLI should interact easily with config files (see [Configuration](#Configuration)), but not _always_
  * expose all config through a `viper` registry
  * leave developers a free-hand on all the knobs and features proposed by `cobra`
2. it should be easier to interact with command line flags of various types
  * simple, declarative registration and binding of flags to config
  * should abstract away the tedious and error prone steps for the registration of flags, binding & defaults
  * allow CLI flags to override this config (12-factors)
  * includes slices, maps and custom flag types (delegated to `github.com/fredbi/gflag`)
3. it should be easier to declare defaults in one single place (for flags, for config)
4. it should be easier to inject external dependencies into the commands tree (config, logger, etc)

### Code style goals

* adopt a functional style (some would say DSL-like), with builder functions for CLI components 
* the tree-like structure of commands should appear visually and obviously in the source code
* remove the need for a typical `init()` to perform all this initialization in the correct order
* remove the need to use package-level variables
* remove the boiler-plate code needed to register, then bind the flags to the config registry
* remove the cognitive burden of remembering all the `GetString()`, `GetBool()` etc methods: `go` now has generics for that
* support the use of generic flags as exposed by [`github.com/fredbi/gflags`](https://github.com/fredbi/gflag), but don't require this
* design with testability in mind: CLI's should be testable with reasonable code coverage

### Ancillary goals

Expose extra packages to address common issues (moving target):
* pre-baked injectable dependencies (logger, ...)
* container sync utilities (file, network port)
* versioning based on go build metadata

### Non-goals

* don't use struct tags: we want to stick to the programmatic approach.
  There are other great libraries around following the struct tags approach
* don't use codegen: we want our code to be readable, not generated


## Configuration

### Desirable features

* load configuration files, using sensible defaults from the powerful `github.com/spf13/viper` package.
* merge configurations, overloading value for a specific environment
* deal with the specifics of merge secrets in config
* help with testing the programs that consume configurations
* leverages all the 12-factor app stuff from `viper`
* leave developers a free-hand if they want to use all the knobs and features proposed by `viper`
  (e.g. dynamic watch, remote config, etc)
* defaults are configurable

### Code style goals

* less boiler plate to deal with `viper` configuration settings and merging

### Non-goals

* Avoid too much of automagically resolving things

  > As much as a like what's available for pythonists with [Dynaconf](https://www.dynaconf.com/),
  > I found myself spending too much time there reading their doc, trying to understand their default settings
  > and figure out whetheir to adopt the default or override.
  >
  > This is a pitfall that is very difficult to avoid, and only experience and feedback will tell.


* Don't want to support older config formats such as `.ini`, `.toml` etc

  > While perfectly doable, I prefer at the moment to focus on having less things to describe and
  > document. So I believe that YAML and JSON are good enough.

* At the moment, no particular goal is set to support secrets via APIs (e.g. Hashicorp's Vault, Azure Vault...)

  > Let's hold on for a bit.
  > At the moment, I am assuming secrets are just plain files (e.g. Kubernetes secret)

