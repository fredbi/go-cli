package wait

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainer(t *testing.T) {
	touchDir, err := os.MkdirTemp("", "touchDir")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll("touchDir")
	})

	t.Run("containerDone should create file", func(t *testing.T) {
		touchFile := filepath.Join(touchDir, "file")

		require.NoError(t, Done(touchFile))

		_, err = os.Stat(touchFile)
		require.NoError(t, err)
	})

	t.Run("containerDone should create dir and file", func(t *testing.T) {
		touchFile := filepath.Join(touchDir, "level", "file")

		require.NoError(t, Done(touchFile))

		_, err = os.Stat(touchFile)
		require.NoError(t, err)
	})
}
