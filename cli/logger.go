package cli

/*
import (
	"github.com/oneconcern/ocpkg/cli/envk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger creates a named zap logger.
//
// It exits the current process upon failure.
func InitLogger(name, level string) *zap.Logger {
	// creates a configured root zap logger

	lc := zap.NewProductionConfig()
	if envk.StringOrDefault("APP_ENV", "local") == "local" {
		lc.Development = true
	}

	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		Die("reading log level config: %v", err)

		return nil
	}

	lc.Level = zap.NewAtomicLevelAt(lvl)
	lc.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	// lc.DisableStacktrace = true

	zlg, err := lc.Build(zap.AddCallerSkip(1))
	if err != nil {
		Die("create logger: %v", err)

		return nil
	}

	zlg = zlg.Named(name)
	zap.ReplaceGlobals(zlg)
	zap.RedirectStdLog(zlg)

	return zlg
}

// LogAndClose initializes a zap logger and provides a trace file (for jaeger support), with a closer for this file.
func LogAndClose(name, level string) (*zap.Logger, func()) {
	zlg := InitLogger(name, level)

	traceCloser, err := traceCfg.InitTraceFile(zlg)
	if err != nil {
		Die("init trace file: %v", err)

		return nil, nil
	}

	return zlg, func() {
		traceCloser()
		_ = zlg.Sync()
	}
}
*/
