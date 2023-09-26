package cli

import (
	"fmt"
	"log"
)

var die = log.Fatalf

// SetDie alters the package level log.Fatalf implementation,
// to be used by Die(sring, ...any).
//
// If fatalFunc is set to nil, calls to Die will issue their message
// with panic instead of log.Fatalf.
//
// This should be used for testing only.
func SetDie(fatalFunc func(string, ...any)) {
	die = fatalFunc
}

// Die exits the current process with some final croak.
// By default, Die is a wrapper around log.Fatalf.
//
// Use SetDie to alter this behavior (e.g. for mocking).
//
// SetDie(nil) will make Die(format, args...) equivalent to
// panic(fmt.Sprintf(format, args...)).
//
// This wraps log.Fatalf, essentially for testing purpose.
func Die(format string, args ...any) {
	if die == nil {
		panic(fmt.Sprintf(format, args...))
	}

	die(format, args...)
}

// Must panic on error
func Must(err error) {
	must(err)
}

func must(err error) {
	if err == nil {
		return
	}

	panic(err)
}
