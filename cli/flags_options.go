package cli

type (
	// FlagOption provides optional behavior to the CLI flags
	// inlined in a Command definition.
	FlagOption func(*flagOption)

	flagOption struct {
		fn         RegisterFunc
		required   bool
		persistent bool
		configKey  string
	}
)

func flagWithOptions(fn RegisterFunc, opts []FlagOption) flagOption {
	fl := flagOption{fn: fn}
	for _, apply := range opts {
		apply(&fl)
	}

	return fl
}

// FlagIsRequired declares the flag as required for the command.
func FlagIsRequired() FlagOption {
	return func(o *flagOption) {
		o.required = true
	}
}

// FlagIsPersistent declares the flag as persistent for the command.
func FlagIsPersistent() FlagOption {
	return func(o *flagOption) {
		o.persistent = true
	}
}

// BindFlagToConfig declares the flag as bound to a configuration key in the viper registry.
func BindFlagToConfig(key string) FlagOption {
	return func(o *flagOption) {
		o.configKey = key
	}
}
