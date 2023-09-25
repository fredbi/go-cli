package injectable

import (
	"context"
	"log/slog"
)

var _ ContextInjectable = &SLogger{}

type SLogger struct {
	logger *slog.Logger
}

func NewSLogger(l *slog.Logger) *SLogger {
	return &SLogger{logger: l}
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

// SLoggerFromContext retrieves a configuration registry from the context.
//
// Optional defaulters can be added to deal with a non-existing config.
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
