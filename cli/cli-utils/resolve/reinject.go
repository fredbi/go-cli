package resolve

import (
	"context"
	"log/slog"

	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ReinjectZapLogger updates the command context with a new logger.
//
// It returns the new context of this command.
//
// This is useful for instance to relevel a logger after CLI flags parsing.
func ReinjectZapLogger(c *cobra.Command, zlg *zap.Logger) context.Context {
	override := injectable.NewZapLogger(zlg)
	c.SetContext(override.Context(c.Context()))

	return c.Context()
}

// ReinjectSLogger updates the command context with a new logger.
//
// It returns the new context of this command.
//
// This is useful for instance to relevel a logger after CLI flags parsing.
func ReinjectSLogger(c *cobra.Command, slg *slog.Logger) context.Context {
	override := injectable.NewSLogger(slg)
	c.SetContext(override.Context(c.Context()))

	return c.Context()
}

// RelevelInjectedZapLogger reinjects into the command's context a relevel logger
// with a new level.
//
// NOTE: the new level can only be higher (more restrictive) than the initially
// configured one, otherwise this is a no-op.
func RelevelInjectedZapLogger(c *cobra.Command, level string, opts ...Option) (context.Context, *zap.Logger) {
	ctx := c.Context()
	o := applyOptions(opts)
	zlg := injectable.ZapLoggerFromContext(ctx, o.zloggerDefaulter)

	if zlg == nil {
		return ctx, nil
	}

	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return ctx, zlg
	}

	zlg = zlg.WithOptions(zap.IncreaseLevel(lvl))

	return ReinjectZapLogger(c, zlg), zlg
}
