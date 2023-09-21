package cli

import (
	"os"

	"github.com/fredbi/go-cli/config"
	"github.com/spf13/viper"
)

var (
	configOptions = []config.Option{}

	// ConfigEnv defines the environment variable used by the Config() function
	// to find the current environment.
	ConfigEnv = "CONFIG_ENV"

	// ConfigDebugEnv defines the environment variable used to instruct the config loader
	// to dump all config keys for debugging.
	ConfigDebugEnv = "DEBUG_CONFIG"
)

// SetConfigOption defines package-level defaults for config options, when using
// ConfigForEnv or Config.
//
// By default, this package doesn't set any particular option and uses all defaults from the config package.
//
// TODO: package-level options should be removed and injected.
func SetConfigOptions(opts ...config.Option) {
	configOptions = opts
}

// ConfigForEnv loads and merge a set of config files for a given environment and applies some default values.
//
// It assumes that config files follow the conventions defined by "github.com/fredbi/go-cli/config".
//
// It dies upon failure.
//
// Environment variable settings:
//   - If the environment variable "DEBUG_CONFIG" is set, the loaded settings are dumped to standard output as JSON.
//   - The environment variable "CONFIG_DIR" defines the folder where the root configuration is located.
func ConfigForEnv(env string, defaulters ...func(*viper.Viper)) *viper.Viper {
	return ConfigForEnvWithOptions(env, configOptions, defaulters...)
}

// ConfigForEnvWithOptions loads and merge a set of config files for a given environment and applies some default values.
//
// This function accepts some config.Options to control where and how the configuration files should be loaded.
func ConfigForEnvWithOptions(env string, opts []config.Option, defaulters ...func(*viper.Viper)) *viper.Viper {
	cfg, err := config.LoaderWithSecrets(opts...).LoadForEnv(env)
	if err != nil {
		die("loading config: %v", err)

		return nil
	}

	for _, defaulter := range defaulters {
		defaulter(cfg)
	}

	if wantsDebug() {
		cfg.Debug()
	}

	return cfg
}

// Config calls ConfigForEnv, with the current environment resolved from the variable "CONFIG_ENV".
func Config(defaulters ...func(*viper.Viper)) *viper.Viper {
	env := config.GetenvOrDefault(ConfigDebugEnv, "")

	return ConfigForEnv(env, defaulters...)
}

func wantsDebug() bool {
	return os.Getenv(ConfigDebugEnv) != ""
}
