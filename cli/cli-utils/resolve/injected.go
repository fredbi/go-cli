package resolve

import (
	"context"
	"log/slog"

	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	// Option to tune the retrieval of injected dependencies from a command's context
	Option func(*options)

	options struct {
		zloggerDefaulter func() *zap.Logger
		sloggerDefaulter func() *slog.Logger
		configDefaulter  func() *viper.Viper
	}
)

func applyOptions(opts []Option) options {
	o := options{}

	for _, apply := range opts {
		apply(&o)
	}

	return o
}

// ResolveInjectZapConfig returns in a single call the zap logger and viper config stored in the
// command context as injectables.
//
// Options may specify defaults whenever these items cannot be found in the context.
func InjectedZapConfig(c *cobra.Command, opts ...Option) (context.Context, *zap.Logger, *viper.Viper) {
	o := applyOptions(opts)
	ctx := c.Context()
	zlg := injectable.ZapLoggerFromContext(ctx, o.zloggerDefaulter)
	cfg := injectable.ConfigFromContext(ctx, o.configDefaulter)

	return ctx, zlg, cfg
}

// ResolveInjectSLogConfig returns in a single call the standard structure logger and viper config stored in the
// command context as injectables.
//
// Options may specify defaults whenever these items cannot be found in the context.
func InjectedSLogConfig(c *cobra.Command, opts ...Option) (context.Context, *slog.Logger, *viper.Viper) {
	o := applyOptions(opts)
	ctx := c.Context()
	slg := injectable.SLoggerFromContext(ctx, o.sloggerDefaulter)
	cfg := injectable.ConfigFromContext(ctx, o.configDefaulter)

	return ctx, slg, cfg
}

// WithZapLoggerDefaulter applies a default logger generation whenever the injected context is not found.
//
// This guarantees the presence of a logger.
func WithZapLoggerDefaulter(defaulter func() *zap.Logger) Option {
	return func(o *options) {
		o.zloggerDefaulter = defaulter
	}
}

// WithSLoggerDefaulter applies a default logger generation whenever the injected context is not found.
//
// This guarantees the presence of a logger.
func WithSLoggerDefaulter(defaulter func() *slog.Logger) Option {
	return func(o *options) {
		o.sloggerDefaulter = defaulter
	}
}

// WithConfigDefaulter applies a default config generation whenever the injected context is not found.
//
// This guarantees on the presence of a config.
func WithConfigDefaulter(defaulter func() *viper.Viper) Option {
	return func(o *options) {
		o.configDefaulter = defaulter
	}
}
