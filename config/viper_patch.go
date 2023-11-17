package config

import "github.com/spf13/viper"

// ViperSub is a patch over viper.Sub(), to solve issue https://github.com/spf13/viper/issues/801.
//
// This works to extract a section as a map, and resolve everything (flags, envs, defaults, etc.).
func ViperSub(cfg *viper.Viper, key string) *viper.Viper {
	if cfg == nil {
		return nil
	}

	configMap := cfg.AllSettings()[key] // we force a resolution of the config here. viper.Sub() doesn't work (any longer?).
	if configMap == nil {
		return nil
	}

	sub := viper.New()

	switch typed := configMap.(type) {
	case map[string]any:
		for k, v := range typed {
			sub.Set(k, v)
		}
		return sub
	default:
		return nil
	}
}
