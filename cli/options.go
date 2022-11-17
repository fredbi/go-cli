package cli

import (
	"github.com/fredbi/go-cli/cli/gflag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
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
		config                *viper.Viper
	}
)

// WithConfig adds a viper.Viper configuration to the command tree.
func WithConfig(cfg *viper.Viper) Option {
	return func(o *options) {
		o.config = cfg
	}
}

// TODO: with logger??

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
		o.flagSetters = append(o.flagSetters, fl)
	}
}

// WithPersistentFlagFunc declares a persistent flag for the command using a RegisterFunc function and some flag options.
func WithPersistentFlagFunc(regFunc RegisterFunc, opts ...FlagOption) Option {
	return func(o *options) {
		fl := flagWithOptions(regFunc, opts)
		o.persistentFlagSetters = append(o.persistentFlagSetters, fl)
	}
}

// TODO: with/without flag address
// WithFlag declares a flag of any type supported by gflag, with some options.
func WithFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name string, defaultValue T, usage string, opts ...FlagOption) Option {
	return WithFlagP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithSliceFlag declares a flag of any slice type supported by gflag, with some options.
func WithSliceFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return WithSliceFlagP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithFlagP declares a flag of any type supported by gflag, with a short name and some options.
func WithFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name, short string, defaultValue T, usage string, opts ...FlagOption) Option {
	return withAnyFlagP(addr, name, short, defaultValue, usage, false, opts...)
}

// WithSliceFlagP declares a flag of any slice type supported by gflag, with a short name and some options.
func WithSliceFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name, short string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return withAnySliceFlagP(addr, name, short, defaultValue, usage, false, opts...)
}

// WithPersistentFlag declares a persistent flag of any type supported by gflag, with some options.
func WithPersistentFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name string, defaultValue T, usage string, opts ...FlagOption) Option {
	return WithPersistentFlagP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithPersistentSliceFlag declares a persistent flag of any slice type supported by gflag, with some options.
func WithPersistentSliceFlag[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return WithPersistentSliceFlagP[T](addr, name, "", defaultValue, usage, opts...)
}

// WithPersistentFlagP declares a persistent flag of any type supported by gflag, with a short name and some options.
func WithPersistentFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name, short string, defaultValue T, usage string, opts ...FlagOption) Option {
	return withAnyFlagP(addr, name, short, defaultValue, usage, true, opts...)
}

// WithPersistentSliceFlagP declares a persistent flag of any type supported by gflag, with a short name and some options.
func WithPersistentSliceFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name, short string, defaultValue []T, usage string, opts ...FlagOption) Option {
	return withAnySliceFlagP(addr, name, short, defaultValue, usage, true, opts...)
}

func withAnyFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *T, name, short string, defaultValue T, usage string, isPersistent bool, opts ...FlagOption) Option {
	regFunc := func(flags *pflag.FlagSet) string {
		if addr == nil {
			addr = new(T)
		}
		v := gflag.NewFlagValue[T](addr, defaultValue)
		fl := flags.VarPF(v, name, short, usage)
		fl.NoOptDefVal = v.NoOptDefVal

		return fl.Name
	}

	return optFuncWithRegister(regFunc, isPersistent, opts)
}

func withAnySliceFlagP[T gflag.FlaggablePrimitives | gflag.FlaggableTypes](addr *[]T, name, short string, defaultValue []T, usage string, isPersistent bool, opts ...FlagOption) Option {
	regFunc := func(flags *pflag.FlagSet) string {
		if addr == nil {
			slice := make([]T, 0, len(defaultValue))
			addr = &slice
		}
		v := gflag.NewFlagSliceValue[T](addr, defaultValue)
		fl := flags.VarPF(v, name, short, usage)

		return fl.Name
	}

	return optFuncWithRegister(regFunc, isPersistent, opts)
}

func optFuncWithRegister(regFunc RegisterFunc, isPersistent bool, opts []FlagOption) Option {
	return func(o *options) {
		fl := flagWithOptions(regFunc, opts)
		if isPersistent {
			o.persistentFlagSetters = append(o.persistentFlagSetters, fl)
		} else {
			o.flagSetters = append(o.flagSetters, fl)
		}
	}
}
