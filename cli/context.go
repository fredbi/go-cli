package cli

import (
	"context"

	"github.com/spf13/viper"
)

type (
	commandCtxKey uint8
)

const (
	ctxConfig commandCtxKey = iota + 1
)

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

// TODO: puts logger in context
