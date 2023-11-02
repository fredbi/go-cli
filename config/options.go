package config

import (
	"io"
	"os"
)

type (
	// Option to configure the config loader.
	Option func(*options)

	options struct {
		basePathFromEnvVar  string
		basePath            string
		envDir              string
		radix               string
		suffix              string
		output              io.Writer
		parentSearchEnabled bool
		skipWatch           bool
	}
)

func defaultOptions(opts []Option) *options {
	o := &options{
		basePathFromEnvVar: "",
		basePath:           DefaultBasePath,
		envDir:             DefaultEnvPath,
		radix:              DefaultConfigRadix,
		suffix:             "",
		skipWatch:          false,
		output:             os.Stdout,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithBasePathFromEnv specifies the environment variable to be used
// to define the search path for configs (default: none).
func WithBasePathFromEnvVar(variable string) Option {
	return func(o *options) {
		o.basePathFromEnvVar = variable
	}
}

// WithBasePath defines the folder of the root configuration file.
//
// Defaults to "." (this default may be altered by changing DefaultBasePath).
func WithBasePath(pth string) Option {
	return func(o *options) {
		o.basePath = pth
	}
}

// WithEnvDir defines the path to environment-specific configs, relative to the base path.
//
// Defaults to "config.d" (this default may be altered by changing DefaultEnvPath).
func WithEnvDir(dir string) Option {
	return func(o *options) {
		o.envDir = dir
	}
}

// WithRadix defines the radix (base name) of a config file.
//
// Default to "config", meaning that the recognized files are: "config.yaml", "config.yml", "config.json" for root files,
// or forms such as "config.xxx.{json|yaml|yml} for env-specific files.
//
// The default may be altered by changing DefaultConfigRadix.
func WithRadix(radix string) Option {
	return func(o *options) {
		o.radix = radix
	}
}

// WithSuffix defines an optional suffix extension to be recognized.
//
// The loader will primarily search for files with the suffix extension, then fall back to no suffix.
//
// Example:
//
//	WithSuffix("dec") will look first for "secrets.yaml.dec",
//	then if none is found, will search for "secrets.yaml".
//
// The default is empty.
func WithSuffix(suffix string) Option {
	return func(o *options) {
		o.suffix = suffix
	}
}

// WithMute discards any logging occurring inside the config loader.
//
// The default is false.
//
// This option is equivalent to WithOutput({io.Discard|os.Stdout}).
func WithMute(mute bool) Option {
	return func(o *options) {
		if mute {
			o.output = io.Discard
		} else {
			o.output = os.Stdout
		}
	}
}

// WithWatch enables the configuration watch.
//
// The default is true.
func WithWatch(canWatch bool) Option {
	return func(o *options) {
		o.skipWatch = !canWatch
	}
}

// WithOutput specifies a io.Writer to output logs during config search.
//
// The default is os.Stdout.
func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

// WithSearchParentDir enables the search for config files in the folder tree
// that contains the current working directory.
//
// This is primarily designed to support test programs loading config files from a source repository.
func WithSearchParentDir(enabled bool) Option {
	return func(o *options) {
		o.parentSearchEnabled = enabled
	}
}
