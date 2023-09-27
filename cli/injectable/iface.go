package injectable

import "context"

type (
	// ContextInjectable knows how to retrieve a value from a Context.
	//
	// Users of the github.com/fredbi/go-cli/cli package can define their own injections via the context.
	//
	// NOTE: every such type should declare their own key type to avoid conflicts inside the context.
	//
	// For example:
	//   type commandCtxKey uint8
	//   const ctxConfig commandCtxKey = iota + 1
	ContextInjectable interface {
		// Context builds a context with the injected value
		Context(context.Context) context.Context

		// FromContext retrieves the injected value from the context.
		//
		// It should work even if the receiver is zero value.
		FromContext(context.Context) interface{}
	}
)
