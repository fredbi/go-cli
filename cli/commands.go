package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/fredbi/go-cli/cli/injectable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	// Command wraps a *cobra.Command with some options to register
	// and bind flags with a functional style.
	Command struct {
		*cobra.Command
		*options
	}
)

// NewCommand builds a new Command wrapping a *cobra.Command.
func NewCommand(cmd *cobra.Command, opts ...Option) *Command {
	options := new(options)
	for _, apply := range opts {
		apply(options)
	}

	c := &Command{
		Command: cmd,
		options: options,
	}

	c.SetContext(c.injectedContext(context.Background()))

	c.pushChildren()
	c.registerFlags()

	// config is a special injected dependency because we can bind it with CLI flags
	if c.config != nil {
		c.bindFlagsWithConfig(c.config)
	}

	if c.withVersion != nil {
		c.Version = c.withVersion()
	}

	for _, apply := range c.cobraOpts {
		if apply == nil {
			continue
		}

		apply(c.Command)
	}

	return c
}

// pushChildren pushes subcommands and propagates the config in their context.
func (c *Command) pushChildren() {
	for _, sub := range c.subs {
		if c.config != nil {
			sub.config = c.config
			sub.SetContext(c.Context())
		}

		for _, apply := range c.cobraOpts {
			if apply == nil {
				continue
			}

			apply(sub.Command)
		}
		c.Command.AddCommand(sub.Command)
	}
}

func (c *Command) injectedContext(ctx context.Context) context.Context {
	if c.config != nil {
		injected := injectable.NewConfig(c.config)

		ctx = injected.Context(ctx)
	}

	for _, injected := range c.injectables {
		ctx = injected.Context(ctx)
	}

	return ctx
}

// Config returns the viper config registry shared by the command tree.
func (c *Command) Config() *viper.Viper {
	return c.config
}

// String provides a short text representation of the command subtree.
func (c *Command) String() string {
	return c.padString(0)
}

func (c *Command) padString(pad int) string {
	var rep strings.Builder
	fmt.Fprintf(&rep, "%s%s",
		strings.Repeat(" ", pad),
		c.Name(),
	)

	subs := c.Commands()
	if len(subs) > 0 {
		fmt.Fprintf(&rep, "\n")

		for _, sub := range subs[:len(subs)-1] {
			fmt.Fprintln(&rep, sub.padString(pad+1))
		}

		fmt.Fprint(&rep, subs[len(subs)-1].padString(pad+1))
	}

	return rep.String()
}

// AddCommand adds child command(s).
func (c *Command) AddCommand(subs ...*Command) {
	c.subs = append(c.subs, subs...)
	c.pushChildren()
}

// Commands returns the child commands.
func (c *Command) Commands() []*Command {
	return c.subs
}

// register command-line flags for this command.
func (c *Command) registerFlags() {
	for _, setter := range c.persistentFlagSetters {
		name := setter.fn(c.PersistentFlags())
		if setter.required {
			must(c.MarkPersistentFlagRequired(name))
		}
		if setter.configKey != "" {
			c.persistentFlagsToBind = append(c.persistentFlagsToBind, binding{name: name, key: setter.configKey})
		}
	}

	for _, setter := range c.flagSetters {
		name := setter.fn(c.Flags())
		if setter.required {
			must(c.MarkFlagRequired(name))
		}
		if setter.configKey != "" {
			c.flagsToBind = append(c.flagsToBind, binding{name: name, key: setter.configKey})
		}
	}
}

// bindFlagsWithConfig binds the command-line flags marked as such to the configuration registry.
//
// This applies recursively to all sub-commands.
//
// It doesn't perform anything if no flags or no config are set for the command
// (use the options WithFlagVar() and WithConfig())
func (c *Command) bindFlagsWithConfig(cfg *viper.Viper) {
	for _, subCommand := range c.Commands() {
		subCommand.bindFlagsWithConfig(cfg)
	}

	for _, bd := range c.persistentFlagsToBind {
		mustBindFromFlagSet(cfg, bd.key, bd.name, c.PersistentFlags())
	}

	for _, bd := range c.flagsToBind {
		mustBindFromFlagSet(cfg, bd.key, bd.name, c.Flags())
	}
}

// Execute the command, with a default context.Background().
//
// It ensures that all injectables are in the context.
func (c *Command) Execute() error {
	return c.ExecuteContext(context.Background())
}

// ExecuteWithArgs is a convenience wrapper to execute a command with preset args.
//
// This is primarily intended for testing commands.
func (c *Command) ExecuteWithArgs(args ...string) error {
	if len(args) > 0 {
		c.SetArgs(args)
	}

	return c.Execute()
}

// ExecuteContext wraps cobra.Command.ExecuteContext().
//
// It ensures that all injectables are in the context.
func (c *Command) ExecuteContext(ctx context.Context) error {
	ctx = c.injectedContext(ctx)

	return c.Command.ExecuteContext(ctx)
}
