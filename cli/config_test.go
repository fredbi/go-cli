package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fredbi/go-cli/config"
	"github.com/stretchr/testify/require"
)

func init() {
	// disable fatal with a global mock
	SetDie(newFatalMock().Fatalf)
}

func TestConfig(t *testing.T) {
	t.Parallel()

	os.Setenv(ConfigDebugEnv, "1")
	t.Cleanup(func() {
		os.Setenv(ConfigDebugEnv, "")
	})

	t.Run("should load config in debug mode", func(t *testing.T) {

		cfg := ConfigForEnvWithOptions("dev", testConfigOptions(t))
		require.NotNil(t, cfg)
	})

	t.Run("should error on load config", func(t *testing.T) {
		cfg := ConfigForEnvWithOptions("dev", append(testConfigOptions(t), config.WithRadix("invalid")))
		require.Nil(t, cfg)
	})
}

func testConfigOptions(t testing.TB) []config.Option {
	// set some config options for testing.
	here, err := os.Getwd()
	require.NoError(t, err)

	fixtures := filepath.Join(here, "fixtures")

	return []config.Option{
		config.WithMute(false),
		config.WithWatch(false),
		config.WithBasePath(fixtures),
	}
}
