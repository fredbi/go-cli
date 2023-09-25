package main

import (
	"fmt"

	"github.com/fredbi/go-cli/cli/version"
)

func main() {
	fmt.Println(version.Resolve())
}
