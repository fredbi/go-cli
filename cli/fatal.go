package cli

import (
	"log"
)

var die = log.Fatalf

// SetDie alters the package level log.Fatalf implementation.
//
// This should be used for testing only.
func SetDie(fatalFunc func(string, ...any)) {
	die = fatalFunc
}

// Die exits the current process with some final croak.
//
// This wraps log.Fatal for convenient testing.
func Die(format string, args ...any) {
	die(format, args...)
}

// Must panic on error.
func Must(err error) {
	must(err)
}

func must(err error) {
	if err == nil {
		return
	}

	panic(err)
}
