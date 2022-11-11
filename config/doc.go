// Package config exposes an opinionated loader for configuration files.
//
// This package exposes a Loadable interface that knows how to load and merge configuration files.
//
// The resulting configuration is made available to consuming applications as a *viper.Viper
// configuration registry.
//
// Supported format: YAML, JSON
//
// Supported file extensions: "yml", "yaml", "json"
//
// Folder structure:
//
//	{base path}/config.yaml                         # <- root configuration
//	{base path}/config.d/{environment}/config.yaml  # <- configuration to merge for environment
package config
