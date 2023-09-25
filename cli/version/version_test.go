package version_test

import (
	"log"
	"testing"

	"github.com/fredbi/go-cli/cli/version"
)

func TestResolve(t *testing.T) {
	log.Println(version.Resolve())
}
