package injectable

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

var (
	_ ContextInjectable = &SLogger{}
	_ ContextInjectable = &ZapLogger{}
)

type (
	SLogger struct {
		logger *slog.Logger
	}

	ZapLogger struct {
		logger *zap.Logger
	}
)

func NewSLogger(l *slog.Logger) *SLogger {
	return &SLogger{logger: l}
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: l}
}

func (c *SLogger) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSLogger, c.logger)
}

func (c SLogger) FromContext(ctx context.Context) interface{} {
	logger, ok := ctx.Value(ctxSLogger).(*slog.Logger)
	if !ok {
		return nil
	}

	return logger
}

func (c *ZapLogger) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSLogger, c.logger)
}

func (c ZapLogger) FromContext(ctx context.Context) interface{} {
	logger, ok := ctx.Value(ctxSLogger).(*slog.Logger)
	if !ok {
		return nil
	}

	return logger
}

// SLoggerFromContext retrieves a standard structured logger from the context.
//
// Optional defaulters can be added to deal with a non-existing logger.
func SLoggerFromContext(ctx context.Context, defaulters ...func() *slog.Logger) *slog.Logger {
	var c SLogger

	logger := c.FromContext(ctx)
	if logger == nil {
		for _, defaulter := range defaulters {
			logger = defaulter()
		}
	}

	return logger.(*slog.Logger)
}

// ZapLoggerFromContext retrieves a zap logger from the context.
//
// Optional defaulters can be added to deal with a non-existing logger.
func ZapLoggerFromContext(ctx context.Context, defaulters ...func() *slog.Logger) *zap.Logger {
	var c ZapLogger

	logger := c.FromContext(ctx)
	if logger == nil {
		for _, defaulter := range defaulters {
			logger = defaulter()
		}
	}

	return logger.(*zap.Logger)
}
