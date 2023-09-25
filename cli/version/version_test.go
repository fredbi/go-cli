package version_test

import (
	"log"
	"testing"

	"github.com/fredbi/go-cli/cli/version"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	require.NotPanics(t, func() {
		log.Println(version.Resolve())
	})
}
