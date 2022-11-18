package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadForTest(t *testing.T) {
	here, err := os.Getwd()
	require.NoError(t, err)

	fixtures := filepath.Join(here, "examples")

	_, err = LoadForTest("dev", WithBasePath(fixtures))
	require.NoError(t, err)
}
