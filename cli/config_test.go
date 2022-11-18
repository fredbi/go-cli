package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fredbi/go-cli/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	onceFatalMock.Do(func() {
		dieMock = newFatalMock()
	})
	dieMock.Register(t)

	t.Cleanup(func() {
		dieMock.Reset()
	})

	os.Setenv(ConfigDebugEnv, "1")
	t.Cleanup(func() {
		os.Setenv(ConfigDebugEnv, "")
	})

	t.Run("should die", func(t *testing.T) {
		t.Parallel()
		dieMock.Register(t)

		Die("test err")
		require.True(t, dieMock.Called(t))
	})

	t.Run("should load config in debug mode", func(t *testing.T) {
		t.Parallel()
		dieMock.Register(t)

		cfg := ConfigForEnvWithOptions("dev", testConfigOptions(t),
			func(v *viper.Viper) {
				v.SetDefault("default.level", "debug")
			},
		)
		require.NotNil(t, cfg)
		require.False(t, dieMock.Called(t))
		require.Equal(t, "debug", cfg.GetString("default.level"))
	})

	t.Run("should error on load config", func(t *testing.T) {
		t.Parallel()
		dieMock.Register(t)

		cfg := ConfigForEnvWithOptions("dev", append(testConfigOptions(t), config.WithRadix("invalid")))
		require.Nil(t, cfg)
		require.True(t, dieMock.Called(t))
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
