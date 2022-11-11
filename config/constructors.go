package config

var (
	// Default values used by this package (may be overridden when loading this package).

	// DefaultPathEnv is the default environment variable used
	// to configure the search path for configurations.
	DefaultPathEnv = "CONFIG_DIR"

	// DefaultConfigRadix is the default basename for a configuration file.
	DefaultConfigRadix = "config"

	// DefaultSecretRadix is the default basename for a secrets file.
	DefaultSecretRadix = "secrets"

	// DefaultEnvPath is the default folder for environment-specific configurations, relative to the base path.
	DefaultEnvPath = "config.d"

	// DefaultSecretSuffix is the default optional suffix to look for unencrypted secrets.
	DefaultSecretSuffix = "dec"
)

// DefaultLoader provides a config loader with all defaults.
//
// * The base path from where to search for configuration files is provided by the environment variable "CONFIG_DIR"
// * Considered root config files are of the form "config.{json|yaml|yml}"
// * Environment-specific config files are found in the "config.d" folder.
// * Considered environment-specific files are of the form "config[.*].{json|yaml|yml}".
// * Environment-specific files are located in the {base path}/config.d/{env} folder.
// * Files are watched for changes.
// * Config loading logging goes to os.Stdout.
//
// Options may be provided to override some of these defaults.
func DefaultLoader(opts ...Option) Loadable {
	options := []Option{
		WithBasePathFromEnvVar(DefaultPathEnv),
	}
	options = append(options, opts...)

	return NewLoader(options...)
}

// SecretsLoader works like the default loader, but includes config files of the form "secrets.{env}.{json|yaml|yml}[.dec]".
//
// This allows to work with files containing secrets, conventionally named "secrets[.*].[json|yaml|yml]".
// Secret files may possibly be temporarily decrypted with a ".dec" extension: this allows
// to work on a local testing environment with secret configurations managed by sops.
func SecretsLoader(opts ...Option) Loadable {
	options := []Option{
		WithBasePathFromEnvVar(DefaultPathEnv),
		WithRadix(DefaultSecretRadix),
		WithSuffix(DefaultSecretSuffix),
	}
	options = append(options, opts...)

	return NewLoader(options...)
}

// LoaderWithSecrets combines configuration files and secrets.
func LoaderWithSecrets(opts ...Option) Loadable {
	return NewCombinedLoader(
		DefaultLoader(opts...),
		SecretsLoader(opts...),
	)
}

// LoaderForTest provides a default config loader intended to be used by test programs.
//
// It mutes internal logging to limit unwanted test verbosity.
// Config watch is disabled: this is convenient when loading many configuration files in different test go routines.
//
// The LoaderForTest searches for configuration files in the tree containing the current working directory, meaning that
// any test program in your source tree may load a test config.
func LoaderForTest(opts ...Option) Loadable {
	options := []Option{
		WithMute(true),
		WithWatch(false),
		WithSearchParentDir(true),
	}
	options = append(options, opts...)

	// TODO: return NewLoaderWithSecrets(options...)
	return NewLoader(options...)
}
