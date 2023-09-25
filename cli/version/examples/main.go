//go:ignore
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fredbi/go-cli/cli/version"
)

func main() {
	spew.Dump(version.Resolve())
}
