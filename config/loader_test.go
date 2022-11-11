package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	t.Parallel()

	t.Run("should resolve to empty config", func(t *testing.T) {
		t.Parallel()

		wd, err := os.MkdirTemp("", "")
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.RemoveAll(wd)
		})

		envsPath := filepath.Join("k8s", "envs")
		testEnv := filepath.Join(wd, envsPath, "test")
		require.NoError(t, os.MkdirAll(testEnv, os.ModePerm))
		otherEnv := filepath.Join(wd, envsPath, "other")
		require.NoError(t, os.MkdirAll(otherEnv, os.ModePerm))

		ldr := NewLoader(WithBasePath(wd), WithEnvDir(envsPath), WithRadix("config"), WithOutput(os.Stderr))
		cfg, err := ldr.LoadForEnv("test")
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("with layout for config files", func(t *testing.T) {
		t.Parallel()

		wd, err := os.MkdirTemp("", "")
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.RemoveAll(wd)
		})

		envsPath := filepath.Join("k8s", "envs")
		testEnv := filepath.Join(wd, envsPath, "test")
		require.NoError(t, os.MkdirAll(testEnv, os.ModePerm))
		otherEnv := filepath.Join(wd, envsPath, "other")
		require.NoError(t, os.MkdirAll(otherEnv, os.ModePerm))

		require.NoError(t, os.MkdirAll(testEnv, os.ModePerm))
		require.NoError(t, os.MkdirAll(otherEnv, os.ModePerm))

		ldr := NewLoader(WithBasePath(wd), WithEnvDir(envsPath), WithRadix("config"), WithMute(true), WithWatch(false))

		t.Run("should load YAML config (root)", func(t *testing.T) {
			const sampleYAML = `
log:
  level: info
metrics:
  opentracing:
    enabled: true
`
			require.NoError(t,
				// root config
				os.WriteFile(filepath.Join(wd, "config.yaml"), []byte(sampleYAML), 0600),
			)

			cfg, err := ldr.LoadForEnv("test")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, "info", cfg.GetString("log.level"))
			assert.True(t, cfg.GetBool("metrics.opentracing.enabled"))
		})

		t.Run("should load JSON config (env: test)", func(t *testing.T) {
			const sampleJSON = `{
  "log": {
    "level": "debug"
    },
  "grpc": {
    "log": {
      "level": "debug",
      "verbosity": 10
    }
  }
}`

			require.NoError(t,
				// config for env "test"
				os.WriteFile(filepath.Join(testEnv, "config.json"), []byte(sampleJSON), 0600),
			)

			cfg, err := ldr.LoadForEnv("test")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, "debug", cfg.GetString("log.level"))
			assert.Equal(t, "debug", cfg.GetString("grpc.log.level"))
			assert.Equal(t, 10, cfg.GetInt("grpc.log.verbosity"))
			assert.True(t, cfg.GetBool("metrics.opentracing.enabled"))
		})

		t.Run("should load env-specific YAML config (env: other)", func(t *testing.T) {
			const envYAML = `
metrics:
  apollotracing:
    enabled: true
`

			require.NoError(t,
				// config for env "other"
				os.WriteFile(filepath.Join(otherEnv, "config.default.yaml"), []byte(envYAML), 0600),
			)

			cfg, err := ldr.LoadForEnv("other")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, "info", cfg.GetString("log.level"))
			assert.Equal(t, "", cfg.GetString("grpc.log.level"))
			assert.True(t, cfg.GetBool("metrics.opentracing.enabled"))
			assert.True(t, cfg.GetBool("metrics.apollotracing.enabled"))
		})

		t.Run("with secrets", func(t *testing.T) {
			t.Run("should load env-specific YAML secret (env: other)", func(t *testing.T) {
				//nolint:gosec
				const secretYAML = `
secrets:
  key: "abc"
metrics:
  apollotracing:
    enabled: false
`
				require.NoError(t,
					// secrets for env "other"
					os.WriteFile(filepath.Join(otherEnv, "secrets.other.yaml.dec"), []byte(secretYAML), 0600),
				)

				ldr := SecretsLoader(WithBasePath(wd), WithEnvDir(envsPath), WithMute(true))
				cfg, err := ldr.LoadForEnv("other")
				require.NoError(t, err)
				require.NotNil(t, cfg)

				assert.Equal(t, "abc", cfg.GetString("secrets.key"))
				assert.Equalf(t, "", cfg.GetString("log.level"),
					"expected secrets loader not to merge root config",
				)
				assert.Falsef(t, cfg.GetBool("metrics.apollotracing.enabled"),
					"expected value in secret to override config value",
				)

				t.Run("should merge config with env-specific YAML secret (env: other)", func(t *testing.T) {
					mdr := LoaderWithSecrets(WithBasePath(wd), WithEnvDir(envsPath), WithMute(true), WithWatch(false))
					mcfg, erm := mdr.LoadForEnv("other")
					require.NoError(t, erm)
					require.NotNil(t, mcfg)

					assert.Equal(t, "abc", mcfg.GetString("secrets.key"))
					assert.Equalf(t, "info", mcfg.GetString("log.level"),
						"expected secrets loader to be merged from root config",
					)
					assert.Falsef(t, mcfg.GetBool("metrics.apollotracing.enabled"),
						"expected value in secret to override config value",
					)
				})

				t.Run("should get base path from environment", func(t *testing.T) {
					mdr := LoaderWithSecrets(WithBasePathFromEnvVar("TEST_DIR"), WithEnvDir(envsPath), WithWatch(false))
					os.Setenv("TEST_DIR", wd)
					mcfg, erm := mdr.LoadForEnv("other")
					require.NoError(t, erm)
					require.NotNil(t, mcfg)

					assert.Equal(t, "abc", mcfg.GetString("secrets.key"))
					assert.Equalf(t, "info", mcfg.GetString("log.level"),
						"expected secrets loader to be merged from root config",
					)
					assert.Falsef(t, mcfg.GetBool("metrics.apollotracing.enabled"),
						"expected value in secret to override config value",
					)
				})
			})

			t.Run("should load env-specific JSON secret (env: test)", func(t *testing.T) {
				//nolint:gosec
				const secretJSON = `{
  "secrets": {
    "key": "xyz"
  }
}
`

				require.NoError(t,
					// secrets for env "test"
					os.WriteFile(filepath.Join(testEnv, "secrets.test.json.dec"), []byte(secretJSON), 0600),
				)

				ldr := SecretsLoader(WithBasePath(wd), WithEnvDir(envsPath))
				cfg, err := ldr.LoadForEnv("test")
				require.NoError(t, err)
				require.NotNil(t, cfg)

				assert.Equal(t, "xyz", cfg.GetString("secrets.key"))
				assert.Equalf(t, "", cfg.GetString("log.level"),
					"expected secrets loader not to merge root config",
				)
			})
		})
	})

	t.Run("with parent layout for config files", func(t *testing.T) {
		t.Parallel()

		wd, err := os.MkdirTemp("", "")
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.RemoveAll(wd)
		})

		deepDir := filepath.Join(wd, "pkg", "sub", "subsub")
		require.NoError(t, os.MkdirAll(deepDir, os.ModePerm))
		envsPath := filepath.Join("k8s", "envs")
		testEnv := filepath.Join(wd, envsPath, "test")
		require.NoError(t, os.MkdirAll(testEnv, os.ModePerm))

		const (
			rootYAML = `
app:
  parallel:
    max: 10
log:
  level: debug
`

			envYAML = `
app:
  parallel:
    max: 11
`
		)

		require.NoError(t,
			// root config
			os.WriteFile(filepath.Join(wd, "config.yaml"), []byte(rootYAML), 0600),
		)
		require.NoError(t,
			// config for env "test"
			os.WriteFile(filepath.Join(testEnv, "config.yaml"), []byte(envYAML), 0600),
		)

		here, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(deepDir))
		defer func() {
			_ = os.Chdir(here)
		}()

		ldr := LoaderForTest(WithEnvDir(envsPath), WithMute(false))

		t.Run("should merge env config", func(t *testing.T) {
			cfg, err := ldr.LoadForEnv("test")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, 11, cfg.GetInt("app.parallel.max"))
			assert.Equal(t, "debug", cfg.GetString("log.level"))
		})

		t.Run("should merge env config when no env is provided", func(t *testing.T) {
			cfg, err := ldr.LoadForEnv("")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, 11, cfg.GetInt("app.parallel.max"))
			assert.Equal(t, "debug", cfg.GetString("log.level"))
		})

		t.Run("should pick root config only", func(t *testing.T) {
			cfg, err := ldr.LoadForEnv("none")
			require.NoError(t, err)
			require.NotNil(t, cfg)

			assert.Equal(t, 10, cfg.GetInt("app.parallel.max"))
			assert.Equal(t, "debug", cfg.GetString("log.level"))
		})
	})
}
