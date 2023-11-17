# cli

A wrapper on top of `cobra.Command` to build
CLIs with a functional style.

## Command builder

* no globals, only functions
* inject `viper.Viper` configuration registry
* allow dependency injection (e.g. logger, etc)

## Utilities

* injectable: logger, config that can be passed via context
* wait: command sync utilities, e.g. to run as deployed containers & sync on events
* version: retrieve version from go module information, or values baked at build-time

## [TODOs](./TODO.md)
