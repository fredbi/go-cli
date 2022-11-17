package containers

import (
	"fmt"
	"os"
	"path/filepath"
)

// ContainerDone signals to other containers that we are done (e.g. proxy container),
// using a touch file to communicate between the containers on the same pod.
//
// This does nothing if the touch file argument is empty.
func ContainerDone(touch string) error {
	if touch == "" {
		return nil
	}

	if dir := filepath.Dir(touch); dir != "." {
		if err := os.MkdirAll(filepath.Dir(touch), 0755); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warn: could not create dir for signal file: %q: %v\n", touch, err)

			return err
		}
	}

	f, err := os.Create(touch)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warn: could not create signal file: %q: %v\n", touch, err)

		return err
	}

	return f.Close()
}
