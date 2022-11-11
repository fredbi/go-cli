package config

import "github.com/spf13/viper"

type (
	// Loadable knows how to load a configuration for an environment.
	Loadable interface {
		LoadForEnv(env string) (*viper.Viper, error)
	}
)

// Load configuration files. The merged outcome is available in a *viper.Viper registry.
func Load(env string, opts ...Option) (*viper.Viper, error) {
	ldr := DefaultLoader(opts...)

	return ldr.LoadForEnv(env)
}

// LoadWithSecrets combines configuration files from a default loader with config from (unencrypted) secrets.
//
// This is a convenient wrapper around LoadForEnv for a configuration composed both regular and secrets files.
//
// # Just like
//
// Usage:
//
//	v, err := config.LoadWithSecrets("dev")
func LoadWithSecrets(env string, opts ...Option) (*viper.Viper, error) {
	ldr := LoaderWithSecrets(opts...)

	return ldr.LoadForEnv(env)
}

// LoadForTest combines configs from a test loader.
//
// This is a convenient wrapper around LoadForEnv for a configuration composed both regular and secrets files.
//
// Typically, test programs would get a configuration from some "test" environment like so:
//
// Usage:
//
//	v, err := config.LoadForTest("test")
func LoadForTest(env string, opts ...Option) (*viper.Viper, error) {
	testOptions := []Option{
		WithRadix(DefaultSecretRadix),
		WithSuffix(DefaultSecretSuffix),
	}
	testOptions = append(testOptions, opts...)

	ldr := NewCombinedLoader(
		LoaderForTest(opts...),
		LoaderForTest(testOptions...),
	)

	return ldr.LoadForEnv(env)
}

// GetenvOrDefault wraps os.Getenv, applying a default value if the environment variable "key" is undefined or empty.
func GetenvOrDefault(key, defaultValue string) string {
	return getenvOrDefault(key, defaultValue)
}
