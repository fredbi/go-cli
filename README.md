# go-cli

A few utilities to build CLI and load config files on top of cobra and viper.

This is based on the seminal work from @casualjim. I am grateful for his inspiring code.

## CLI

### Typical start CLI of a kubernetes deployment

## Configuration

This package proposes an opinionated approach to config files, with loaders that know about environments and secrets.

### Approach to configuration

The goal of this approach is to merge configuration files with environment specifics.

Applications will then be able to consume the settings from a viper configuration registry.

Supported format: YAML, JSON

Support file extensions: "yml", "yaml", "json"

Folder structure:
```
{base path}/config.yaml  # <- root configuration
{base path}/config.d/{environment}/config.yaml  # <- configuration to merge for environment
```

When using default settings (these are configurable), the base path is defined by the `CONFIG_DIR` environment variable.

### Features

* sensible defaults for minimul boiler plate
* most defaults are configurable
* extensible
* the viper object may be watched dynamically
* can merge plain config with secrets
* helper method to easily inject configurations in test programs

### Typical configuration of a kubernetes deployment

Typically, the configuration files are held in one or several `Configmap` resources, mounted by your deployed container.
Most likely, secret files will be mounted from `Secret` resources.

(... to be continued...)

## TODOs

* [ ] Better testability for local secrets
* [ ] Support for viper remote configurations (e.g. consul, etc..)
* [ ] k8s examples with config maps

### Side notes

* Why several modules?
  I wanted to build a single repo with a few loosely related utilities. 
  I realized that putting this together resulted in a large set of dependencies,
  which is almost always annoying to result when you consume only a few of the exposed packages.
  By declaring smaller `go.mod` manifests, I've somewhat reduced this burden on consuming packages.
