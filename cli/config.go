package cli

import (
	"context"
	"os"

	"github.com/fredbi/go-cli/config"
	// jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
)

type (
	commandCtxKey uint8
)

const (
	ctxConfig commandCtxKey = iota + 1
)

var (
	configOptions = []config.Option{}
	// json          jsoniter.API
	debug bool

	// ConfigEnv defines the environment variable used by the Config() function
	// to find the current environment.
	ConfigEnv = "CONFIG_ENV"

	// ConfigDebugEnv defines the environment variable used to instruct the config loader
	// to dump all config keys for debugging.
	ConfigDebugEnv = "DEBUG_CONFIG"
)

/*
func init() {
	json = jsoniter.ConfigFastest
}
*/

// SetConfigOption defines package-level defaults for config options, when using
// ConfigForEnv or Config.
//
// By default, this package sets no option and uses all defaults from the config package.
func SetConfigOptions(opts ...config.Option) {
	configOptions = opts
}

// ConfigForEnv loads and merge a set of config files for a given environment and applies some defaults,
//
// It assumes that config files follow the conventions defined by "github.com/fredbi/go-cli/config".
//
// It dies upon failure.
//
// Environment variable settings
//   - If the environment variable "DEBUG_CONFIG" is set, the loaded settings are dumped to standard output as JSON.
//   - The environment variable "CONFIG_DIR" defines the folder where the root configuration is located.
func ConfigForEnv(env string, defaulters ...func(*viper.Viper)) *viper.Viper {
	cfg, err := config.LoaderWithSecrets(configOptions...).LoadForEnv(env)
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

// ContextWithConfig puts a configuration registry in the context.
func ContextWithConfig(ctx context.Context, cfg *viper.Viper) context.Context {
	return context.WithValue(ctx, ctxConfig, cfg)
}

// ConfigFromContext retrieves a configuration registry from the context.
func ConfigFromContext(ctx context.Context) *viper.Viper {
	cfg, ok := ctx.Value(ctxConfig).(*viper.Viper)
	if !ok {
		return nil
	}

	return cfg
}

func wantsDebug() bool {
	return os.Getenv(ConfigDebugEnv) != "" || debug
}
