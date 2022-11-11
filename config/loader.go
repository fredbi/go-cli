package config

import (
	"fmt"
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type (
	// Loader loads and merge configuration files.
	Loader struct {
		*options
	}

	// CombinedLoader loads and merge configuration files using a collection of loaders.
	CombinedLoader struct {
		loaders []Loadable
	}
)

var (
	_ Loadable = &Loader{}
	_ Loadable = &CombinedLoader{}

	json jsoniter.API
)

func init() {
	json = jsoniter.ConfigFastest
}

// NewLoader creates a new loader for config files.
func NewLoader(opts ...Option) *Loader {
	return &Loader{
		options: defaultOptions(opts),
	}
}

func (l *Loader) LoadForEnv(env string) (*viper.Viper, error) {
	defaultCfg := viper.New()
	defaultCfg.SetConfigName(l.radix)

	if l.parentSearchEnabled {
		// override path configuration and look in the tree containing the current working directory.
		//
		// This should be reserved to testing programs loading configs from a source repository.
		base, err := findParentDir(l.radix, cfgTypes())
		if err != nil {
			return nil, err
		}

		l.basePath = base
	}

	if key := l.basePathFromEnvVar; key != "" {
		l.basePath = getenvOrDefault(key, l.basePath)
	}
	defaultCfg.AddConfigPath(l.basePath)

	// load the root config file
	canWatch, err := l.loadConfig(defaultCfg)
	if err != nil {
		return nil, err
	}

	// load for environment-specific files
	search := filepath.Join(l.basePath, l.envDir, env)
	fi, err := os.Stat(search)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if fi != nil && fi.IsDir() {
		err = filepath.Walk(search, l.walkFunc(defaultCfg))
		if err != nil {
			return nil, err
		}
	}

	defaultCfg.AutomaticEnv()
	if canWatch && !l.skipWatch {
		defaultCfg.WatchConfig()
	}

	return defaultCfg, nil
}

func (l *Loader) walkFunc(cfg *viper.Viper) filepath.WalkFunc {
	return func(pth string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, ext := range cfgTypes() {
			if !matches(l.radix, mkExt(ext, l.suffix), info.Name()) {
				if !matches(l.radix, mkExt(ext, ""), info.Name()) {
					continue
				}
			}

			toMerge, err := l.parseConfigFromExt(pth, ext)
			if err != nil {
				return err
			}

			if err = cfg.MergeConfigMap(toMerge); err != nil {
				return err
			}
		}

		return nil
	}
}

func (l *Loader) loadConfig(cfg *viper.Viper) (bool, error) {
	err := cfg.ReadInConfig()
	if err == nil {
		_, _ = fmt.Fprintln(l.output, "using config file:", cfg.ConfigFileUsed())

		return true, nil
	}

	switch err.(type) {
	case viper.ConfigFileNotFoundError, *viper.ConfigFileNotFoundError:
		_, _ = fmt.Fprintf(l.output,
			"warn: no config file found (pattern: %q), defaulting to empty config\n",
			l.radix,
		)

		return false, nil

	default:
		return false, err
	}
}

func (l *Loader) parseConfigFromExt(pth, ext string) (map[string]interface{}, error) {
	toMerge := make(map[string]interface{})
	_, _ = fmt.Fprintln(l.output, "including config file:", pth)

	buf, err := os.ReadFile(pth)
	if err != nil {
		return nil, err
	}

	switch ext {
	case "yaml", "yml":
		if err = yaml.Unmarshal(buf, &toMerge); err != nil {
			return nil, err
		}
	case "json":
		if err = json.Unmarshal(buf, &toMerge); err != nil {
			return nil, err
		}
	}

	return toMerge, nil
}

// NewCombinedLoader builds a compound loader considering several Loadable in the provided order.
func NewCombinedLoader(loaders ...Loadable) *CombinedLoader {
	return &CombinedLoader{loaders: loaders}
}

func (c *CombinedLoader) LoadForEnv(env string) (*viper.Viper, error) {
	var result *viper.Viper

	for _, loader := range c.loaders {
		toMerge, err := loader.LoadForEnv(env)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = toMerge
		} else if err = result.MergeConfigMap(toMerge.AllSettings()); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func matches(name, ext, input string) bool {
	b, err := filepath.Match(name+"."+ext, input)
	if err == nil && b {
		return true
	}

	b, err = filepath.Match(name+".*."+ext, input)
	if err == nil && b {
		return true
	}

	return false
}

func getenvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func cfgTypes() []string {
	return []string{"yaml", "yml", "json"}
}

func mkExt(ext, suffix string) string {
	if suffix == "" {
		return ext
	}

	return ext + "." + suffix
}

// findInParentDir explores the current directory and its parents to look for a root config file.
//
// This is useful for test programs looking for a config in the repo tree.
func findParentDir(configFile string, allowedExts []string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var (
		pth   string
		found bool
	)

LOOP:
	for cwd != "/" {
		for _, ext := range allowedExts {
			search := filepath.Join(cwd, mkExt(configFile, ext))
			_, err = os.Stat(search)
			if err == nil {
				pth = cwd
				found = true

				break LOOP
			}

			cwd = filepath.Dir(cwd)
		}
	}

	if !found {
		return "", fmt.Errorf("cannot find config location (%s) in current parent tree", configFile)
	}

	return pth, err
}
