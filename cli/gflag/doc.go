// Package gflag exposes generic types to deal with flags.
//
// the new types can be readily used as extensions to the github.com/spf13/flag package.
//
// Most of the current pflag functionality can be obtained from gflag using generic types.
//
// There are a few exceptions though:
// * []byte semantics as base64-encoded string is not available yet (will be, as part of the extensions package).
// * int semantics a increment/decrement count is not avaiable directly (it is available in the extensions package)
// * map semantics (StringToXXX family of pflag types) are not supported yet.
package gflag
