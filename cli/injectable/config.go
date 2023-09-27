package injectable

import (
	"context"

	"github.com/spf13/viper"
)

type (
	commandCtxKey uint8

	// Config can wrap a viper.Viper configuration in the context
	Config struct {
		config *viper.Viper
	}

	// TODO: puts logger in context
	// ZapLogger can wrap a zap logger in the context
	ZapLogger struct {
	}
)

const (
	ctxConfig commandCtxKey = iota + 1
	ctxSLogger
)

var _ ContextInjectable = &Config{}

func NewConfig(v *viper.Viper) *Config {
	return &Config{config: v}
}

func (c *Config) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxConfig, c.config)
}

func (c Config) FromContext(ctx context.Context) interface{} {
	cfg, ok := ctx.Value(ctxConfig).(*viper.Viper)
	if !ok {
		return nil
	}

	return cfg
}

// ConfigFromContext retrieves a configuration registry from the context.
//
// An optional defaulter can be added to deal with a non-existing config.
func ConfigFromContext(ctx context.Context, defaulters ...func() *viper.Viper) *viper.Viper {
	var c Config

	cfg := c.FromContext(ctx)
	if cfg == nil {
		for _, defaulter := range defaulters {
			if defaulter == nil {
				continue
			}

			cfg = defaulter()
		}
	}

	return cfg.(*viper.Viper)
}
