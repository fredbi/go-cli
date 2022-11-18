package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	FlagOption func(*flagOption)
)

// MustBindPlag binds a key in a *viper.Viper registry to a command-line flag (*pflag.Flag).
//
// Dies on error. This happens if the flag is nil.
func MustBindPFlag(cfg *viper.Viper, key string, flag *pflag.Flag) {
	mustBindPFlag(cfg, key, flag)
}

// MustBindFromFlagSet binds a key in a *viper.Viper registry to command-line flag found in a flag set (*pflag.FlagSet).
//
// Dies on error. This happens if the flag set is nil or if the requested flag has not been registered
// in the flag set yet.
func MustBindFromFlagSet(cfg *viper.Viper, key, flagName string, flags *pflag.FlagSet) {
	mustBindFromFlagSet(cfg, key, flagName, flags)
}

func mustBindFromFlagSet(cfg *viper.Viper, key, flagName string, flags *pflag.FlagSet) {
	if flags == nil {
		die("cannot bind on nil flags set")

		return
	}

	flag := flags.Lookup(flagName)
	if flag == nil {
		die("binding unknown pflag to key: %v", key)

		return
	}

	mustBindPFlag(cfg, key, flag)
}

func mustBindPFlag(cfg *viper.Viper, key string, flag *pflag.Flag) {
	if flag == nil {
		die("binding unknown pflag to key: %v", key)

		return
	}

	if err := cfg.BindPFlag(key, flag); err != nil {
		die("binding pflag %s to key %s: %v", flag.Name, key, err)

		return
	}
}
