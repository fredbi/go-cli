package cli

import (
	"encoding/json"

	"github.com/fredbi/gflag"
	"github.com/fredbi/go-cli/cli/cli-utils/version"
	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	// Option configures a Command.
	Option func(*options)

	// RegisterFunc registers a flag to a provided pflag.FlagSet and returns the flag name.
	RegisterFunc func(*pflag.FlagSet) string

	binding struct {
		name string
		key  string
	}

	options struct {
		flagSetters           []flagOption
		persistentFlagSetters []flagOption
		flagsToBind           []binding
		persistentFlagsToBind []binding
		subs                  []*Command
		withVersion           func() string

		// injected dependencies
		config      *viper.Viper
		injectables []injectable.ContextInjectable
	}
)

// WithConfig adds a viper.Viper configuration to the command tree.
//
// The config can be retrieved from the context of the command, using
// injectable.ConfigFromContext().
func WithConfig(cfg *viper.Viper) Option {
	return func(o *options) {
		o.config = cfg
	}
}

// WithVersion wires a "--version" flag from a function evaluated at command construction time.
func WithVersion(versionFunc func() string) Option {
	return func(o *options) {
		o.withVersion = versionFunc
	}
}

// WithAutoVersion wires a "--version" flag using version.Resolve().
//
// This collects versioning from package-level variables set at build time OR
// go module metadata.
func WithAutoVersion() Option {
	return func(o *options) {
		o.withVersion = func() string {
			jazon, _ := json.Marshal(version.Resolve())

			return string(jazon)
		}
	}
}

// WithSubCommands adds child commands.
func WithSubCommands(subs ...*Command) Option {
	return func(o *options) {
		o.subs = append(o.subs, subs...)
	}
}

// WithFlagFunc declares a command flag using a RegisterFunc function and some flag options.
func WithFlagFunc(regFunc RegisterFunc, opts ...FlagOption) Option {
	return func(o *options) {
		fl := flagWithOptions(regFunc, opts)
		if fl.persistent {
			o.persistentFlagSetters = append(o.persistentFlagSetters, fl)
		} else {
			o.flagSetters = append(o.flagSetters, fl)
		}
	}
}

// WithBindFlagsToConfig binds all (non persistent) flags in the provided map to a config key.
//
// The map is like map[string][string]{"flag-name": "viper-key"}
func WithBindFlagsToConfig(bindings map[string]string) Option {
	return func(o *options) {
		for name, configKey := range bindings {
			o.flagsToBind = append(o.flagsToBind, binding{name: name, key: configKey})
		}
	}
}

// WithBindPersistentFlagsToConfig binds all persistent flags in the provided map to a config key.
//
// The map is like map[string][string]{"flag-name": "viper-key"}
func WithBindPersistentFlagsToConfig(bindings map[string]string) Option {
	return func(o *options) {
		for name, configKey := range bindings {
			o.persistentFlagsToBind = append(o.persistentFlagsToBind, binding{name: name, key: configKey})
		}
	}
}

// WithFlag declares a flag of any type supported by gflag, with some options.
//
// The pointer to the flag value is allocated automatically.
func WithFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](name string, defaultValue T, usage string, opts ...FlagOption) Option {
	return WithFlagP[T](name, "", defaultValue, usage, opts...)
}

// WithFlagP declares a flag of any type supported by gflag, with a shorthand name and some options.
//
// The pointer to the flag value is allocated automatically.
func WithFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](name, shorthand string, defaultValue T, usage string, opts ...FlagOption) Option {
	return WithFlagVarP[T](nil, name, shorthand, defaultValue, usage, opts...)
}

// WithFlagVar declares a flag of any type supported by gflag, with some options.
//
// The pointer to the flag value is provided explicitly.
func WithFlagVar[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name string, defaultValue T, usage string, opts ...FlagOption) Option {
	return WithFlagVarP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithFlagVarP declares a flag of any type supported by gflag, with a shorthand name and some options.
//
// The pointer to the flag value is provided explicitly.
func WithFlagVarP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name, shorthand string, defaultValue T, usage string, opts ...FlagOption) Option {
	return withAnyFlagP(addr, name, shorthand, defaultValue, usage, opts...)
}

// WithSliceFlag declares a flag of any slice type supported by gflag, with some options.
//
// The pointer to the flag value is allocated automatically.
func WithSliceFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](name string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return WithSliceFlagP[T](name, "", defaultValue, usage, opts...)
}

// WithSliceFlagP declares a flag of any slice type supported by gflag, with a shorthand name and some options.
//
// The pointer to the flag value is allocated automatically.
func WithSliceFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](name, shorthand string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return WithSliceFlagVarP[T](nil, name, shorthand, defaultValue, usage, opts...)
}

// WithSliceFlagVar declares a flag of any slice type supported by gflag, with some options.
//
// The pointer to the flag value is provided explicitly.
func WithSliceFlagVar[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return WithSliceFlagVarP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithSliceFlagVarP declares a flag of any slice type supported by gflag, with a shorthand name and some options.
func WithSliceFlagVarP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name, shorthand string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return withAnySliceFlagP(addr, name, shorthand, defaultValue, usage, opts...)
}

// WithInjectables adds dependencies to be injected in the context of the command.
//
// For each injectable, its Context() method will be called in the specified order to enrich the context of the command.
//
// NOTE: the config registry is a special dependency because it may bind to CLI flags.
//
// Configuration may be injected directly with the more explicit WithConfig() method.
func WithInjectables(injectables ...injectable.ContextInjectable) Option {
	return func(o *options) {
		o.injectables = append(o.injectables, injectables...)
	}
}

func withAnyFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name, shorthand string, defaultValue T, usage string, opts ...FlagOption) Option {
	regFunc := func(flags *pflag.FlagSet) string {
		if addr == nil {
			addr = new(T)
		}
		v := gflag.NewFlagValue[T](addr, defaultValue)
		fl := flags.VarPF(v, name, shorthand, usage)
		fl.NoOptDefVal = v.NoOptDefVal

		return fl.Name
	}

	return optFuncWithRegister(regFunc, opts)
}

func withAnySliceFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name, short string, defaultValue []T, usage string, opts ...FlagOption) Option {
	regFunc := func(flags *pflag.FlagSet) string {
		if addr == nil {
			slice := make([]T, 0, len(defaultValue))
			addr = &slice
		}
		v := gflag.NewFlagSliceValue[T](addr, defaultValue)
		fl := flags.VarPF(v, name, short, usage)

		return fl.Name
	}

	return optFuncWithRegister(regFunc, opts)
}

func optFuncWithRegister(regFunc RegisterFunc, opts []FlagOption) Option {
	return func(o *options) {
		fl := flagWithOptions(regFunc, opts)
		if fl.persistent {
			o.persistentFlagSetters = append(o.persistentFlagSetters, fl)
		} else {
			o.flagSetters = append(o.flagSetters, fl)
		}
	}
}
